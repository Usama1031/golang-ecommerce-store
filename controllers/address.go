package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/usama1031/golang-ecommerce-store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addresses models.Address

		addresses.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}

		unwindStage := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage})

		if err != nil {
			c.IndentedJSON(500, "Internal server error!")
			return
		}
		defer pointcursor.Close(ctx)

		var addressInfo []bson.M
		if err = pointcursor.All(ctx, &addressInfo); err != nil {

			panic(err)
		}

		var size int32

		for _, address_no := range addressInfo {
			count := address_no["count"]
			size = count.(int32)
		}

		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}

			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

			_, err := UserCollection.UpdateOne(ctx, filter, update)

			if err != nil {
				fmt.Println(err)
				return
			}
			c.IndentedJSON(200, "Address added successfully")
		} else {
			c.IndentedJSON(400, "Only two addressees allowed")
		}
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address

		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.0.house_nmae", Value: editaddress.House}, {Key: "address.0.street_name", Value: editaddress.Street}, {Key: "address.0.city_name", Value: editaddress.City}, {Key: "address.0.pin_code", Value: editaddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(500, "Something went wrong!")
			return
		}

		c.IndentedJSON(200, "Successfully updated home address!")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address

		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address.1.house_nmae", Value: editaddress.House}, {Key: "address.1.street_name", Value: editaddress.Street}, {Key: "address.1.city_name", Value: editaddress.City}, {Key: "address.1.pin_code", Value: editaddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(500, "Something went wrong!")
			return
		}

		c.IndentedJSON(200, "Successfully updated work address!")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}

		update := bson.D{
			primitive.E{
				Key: "$set",
				Value: bson.D{
					primitive.E{Key: "address", Value: addresses},
				},
			},
		}
		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "wrong command")
			return
		}

		c.IndentedJSON(200, "Successfully delete the address!")

	}
}
