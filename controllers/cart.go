package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/usama1031/golang-ecommerce-store/database"
	"github.com/usama1031/golang-ecommerce-store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")

		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			return
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))

			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully added to the cart!")

	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")

		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			return
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))

			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.RemoveItemFromCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully removed item from cart!")

	}
}

func (app *Application) ViewCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Query("id")

		if user_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var filledCart models.User
		err := app.userCollection.FindOne(ctx, bson.M{"_id": usert_id}).Decode(&filledCart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return
		}

		// filter_match := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: usert_id}}}}

		// unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}

		// grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		// pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		// if err != nil {
		// 	log.Println(err)
		// }

		// var listing []bson.M
		// if err = pointcursor.All(ctx, &listing); err != nil {
		// 	log.Println(err)
		// 	c.AbortWithStatus(http.StatusInternalServerError)

		// }

		// for _, json := range listing {
		// 	c.IndentedJSON(200, json["total"])
		// 	c.IndentedJSON(200, filledCart.UserCart)
		// }

		var total float64 = 0

		for _, item := range filledCart.UserCart {
			total += float64(item.Price)
		}

		c.JSON(http.StatusOK, gin.H{
			"user_cart": filledCart.UserCart,
			"total":     total,
		})

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Panic("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))

			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)

		}

		c.IndentedJSON(200, "Successfully placed the order!")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")

		if productQueryID == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			return
		}

		userQueryID := c.Query("userID")

		if userQueryID == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))

			return
		}

		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err != nil {
			log.Println(err)
			log.Println("Error while converting product id")
			return
		}
		log.Println("Product ID", productID)

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = database.InstantBuyer(ctx, app.userCollection, app.prodCollection, productID, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully placed the order!")
	}
}
