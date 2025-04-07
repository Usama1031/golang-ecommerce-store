package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	database "github.com/usama1031/golang-ecommerce-store/database"
	helper "github.com/usama1031/golang-ecommerce-store/helper"
	models "github.com/usama1031/golang-ecommerce-store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "users")

var ProductCollection *mongo.Collection = database.ProductData(database.Client, "products")

var validate = validator.New()

func HashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedBytes)

}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Password is incorrect!")
		check = false
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		filter := bson.M{"$or": []bson.M{{"email": user.Email}, {"phone": user.Phone}}}

		count, err := UserCollection.CountDocuments(ctx, filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking user existence!"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone already exists!"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		var token, refreshToken, _ = helper.GenerateToken(*user.Email, *user.First_name, *user.Last_name, user.User_id)

		user.Token = &token
		user.Refresh_Token = &refreshToken

		user.UserCart = make([]models.ProductUser, 0)

		user.Address_details = make([]models.Address, 0)

		user.Order_status = make([]models.Order, 0)

		_, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User was not created. Please try again!"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully signed-up!")

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User was not found!"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helper.GenerateToken(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		c.SetCookie("token", token, 3600, "/", "", false, true)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusFound, foundUser)
	}
}

func AdminAddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		product.Product_ID = primitive.NewObjectID()

		_, err := ProductCollection.InsertOne(ctx, product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not inserted"})
			return
		}

		c.JSON(http.StatusOK, "Successfully added")

	}
}

func ProductViewer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var productList []models.Product

		cursor, err := ProductCollection.Find(ctx, bson.M{})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Couldn't get the product list!")
			return
		}

		err = cursor.All(ctx, &productList)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		c.IndentedJSON(200, productList)

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var searchProducts []models.Product

		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("Query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			c.Abort()
			return
		}

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam, "$options": "i"}})

		if err != nil {
			c.IndentedJSON(404, "not found, something went wrong")
			return
		}

		err = searchquerydb.All(ctx, &searchProducts)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid")
			return
		}

		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}

		c.IndentedJSON(200, searchProducts)
	}
}
