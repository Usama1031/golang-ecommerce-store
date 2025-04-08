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

		var userAddress models.Address

		userAddress.Address_ID = primitive.NewObjectID()

		if err = c.BindJSON(&userAddress); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		var user models.User

		err = UserCollection.FindOne(ctx, bson.D{{Key: "_id", Value: address}}).Decode(&user)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if len(user.Address_details) > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "User already has an address. Please edit to change your address!"})
			return
		}

		filter := bson.M{"_id": address}

		update := bson.M{
			"$push": bson.M{
				"address_details": userAddress,
			},
		}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, "Error updating address")
			return
		}

		// var addressInfo []bson.M
		// if err = pointcursor.All(ctx, &addressInfo); err != nil {

		// 	panic(err)
		// }

		// var size int32

		// for _, address_no := range addressInfo {
		// 	count := address_no["count"]
		// 	size = count.(int32)
		// }

		// if size < 2 {
		// 	filter := bson.D{primitive.E{Key: "_id", Value: address}}

		// 	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		// 	_, err := UserCollection.UpdateOne(ctx, filter, update)

		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// 	c.IndentedJSON(200, "Address added successfully")
		// } else {
		// 	c.IndentedJSON(400, "Only two addressees allowed")
		// }

		c.IndentedJSON(200, "Address added successfully")
	}
}

func EditAddress() gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while converting the user id"})
		}

		var editaddress models.Address

		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		filter := bson.D{{Key: "_id", Value: usert_id}}

		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "address.0.house_name", Value: editaddress.House},
			{Key: "address.0.street_name", Value: editaddress.Street},
			{Key: "address.0.city_name", Value: editaddress.City},
			{Key: "address.0.pincode", Value: editaddress.Pincode},
		}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}

		c.IndentedJSON(200, "Address added successfully")
	}
}

// func EditWorkAddress() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		user_id := c.Query("id")
// 		if user_id == "" {
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index!"})
// 			c.Abort()
// 			return
// 		}

// 		usert_id, err := primitive.ObjectIDFromHex(user_id)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while converting the user id"})
// 		}

// 		var editaddress models.Address

// 		if err := c.BindJSON(&editaddress); err != nil {
// 			c.IndentedJSON(http.StatusBadRequest, err.Error())
// 		}

// 		filter := bson.D{{Key: "_id", Value: usert_id}}
// 		update := bson.D{{Key: "$set", Value: bson.D{
// 			{Key: "address.1.house_name", Value: editaddress.House},
// 			{Key: "address.1.street_name", Value: editaddress.Street},
// 			{Key: "address.1.city_name", Value: editaddress.City},
// 			{Key: "address.1.pincode", Value: editaddress.Pincode},
// 		}}}

// 		_, err = UserCollection.UpdateOne(ctx, filter, update)

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"error": "Successfully updated the work address"})
// 	}
// }

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
