package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/usama1031/golang-ecommerce-store/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.RegisteredClaims
}

var UserCollection *mongo.Collection = database.UserData(database.Client, "users")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateToken(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
		},
	}

	refreshClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 168)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Fatal(err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Fatal(err)
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updatedObj primitive.D

	updatedObj = append(updatedObj, bson.E{Key: "token", Value: signedToken})
	updatedObj = append(updatedObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true

	filter := bson.M{"user_id": userId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserCollection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updatedObj},
	}, &opt)

	if err != nil {
		log.Panic(err)
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = fmt.Sprintf("The token is invalid!")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		msg = fmt.Sprintf("Token is expired")
		return
	}

	return claims, msg
}
