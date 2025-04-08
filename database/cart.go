package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/usama1031/golang-ecommerce-store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCannotFindProduct        = errors.New("cannot find the product")
	ErrCannotDecodeProducts     = errors.New("cannot find the product")
	ErrUserIdIsNotValid         = errors.New("this user is not valid")
	ErrCannotUpdateUser         = errors.New("cannot add this product to the cart")
	ErrCannotRemoveItemFromCart = errors.New("cannot remove this item from the cart")
	ErrCannotGetItem            = errors.New("unable to get the item from cart")
	ErrCannotBuyCartItem        = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	var product models.ProductUser

	err := prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)

	if err != nil {
		log.Println("Product not found:", err)
		return ErrCannotFindProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println("Invalid user id:", err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "usercart", Value: bson.D{
			{Key: "$each", Value: []models.ProductUser{product}},
		}},
	}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println("Error updating the user cart:", err)
		return ErrCannotUpdateUser
	}

	return nil

}

func RemoveItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		return ErrCannotRemoveItemFromCart
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_at = time.Now()
	orderCart.Order_cart = make([]models.ProductUser, 0)

	orderCart.Payment_method.COD = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$usercart"}}}}

	grouping := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$usercart.price"}}}}}}

	currectRes, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		unwind, grouping,
	})

	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err := currectRes.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var total_price int32

	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	orderCart.Price = int(total_price)

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: orderCart}}}}

	_, err = userCollection.UpdateMany(ctx, filter, update)

	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&getCartItems)

	if err != nil {
		log.Println(err)
	}

	filter1 := bson.D{{Key: "_id", Value: id}}

	update1 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getCartItems.UserCart}}}

	_, err = userCollection.UpdateOne(ctx, filter1, update1)

	if err != nil {
		log.Panic(err)
	}

	usercart_empty := make([]models.ProductUser, 0)

	filter2 := bson.D{{Key: "_id", Value: id}}

	update2 := bson.D{{Key: "$set", Value: bson.D{{Key: "usercart", Value: usercart_empty}}}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)

	if err != nil {
		return ErrCannotBuyCartItem
	}

	return nil
}

// func InstantBuyer(ctx context.Context, userCollection, prodCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

// 	id, err := primitive.ObjectIDFromHex(userID)

// 	if err != nil {
// 		log.Println(err)
// 		return ErrUserIdIsNotValid
// 	}

// 	var product_details models.ProductUser
// 	var order_details models.Order

// 	order_details.Order_ID = primitive.NewObjectID()

// 	order_details.Ordered_at = time.Now()
// 	order_details.Order_cart = make([]models.ProductUser, 0)
// 	order_details.Payment_method.COD = true

// 	err = prodCollection.FindOne(ctx, bson.D{{Key: "_id", Value: productID}}).Decode(&product_details)

// 	if err != nil {
// 		log.Println("Product not found", err)
// 		return ErrCannotFindProduct
// 	}

// 	order_details.Order_cart = append(order_details.Order_cart, product_details)

// 	order_details.Price = product_details.Price

// 	filter := bson.D{{Key: "_id", Value: id}}

// 	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: order_details}}}}

// 	_, err = userCollection.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		log.Println("Error updating user collection", err)
// 		return err
// 	}

// 	filter1 := bson.D{{Key: "_id", Value: id}}

// 	update1 := bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}

// 	_, err = userCollection.UpdateOne(ctx, filter1, update1)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return nil

// }

func InstantBuyer(ctx context.Context, userCollection, prodCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	log.Println("Using prodCollection:", prodCollection.Name())

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	log.Println("Got this product ID from controllers:", productID)

	var product_details models.ProductUser
	err = prodCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product_details)

	if err != nil {
		log.Println("Product not found", err)
		log.Println(product_details)
		return ErrCannotFindProduct
	}

	// Retrieve the user's current cart from the user collection
	var user_details models.User
	err = userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&user_details)
	if err != nil {
		log.Println("User not found", err)
		return ErrUserIdIsNotValid
	}

	// Prepare the order details
	var order_details models.Order
	order_details.Order_ID = primitive.NewObjectID()
	order_details.Ordered_at = time.Now()
	order_details.Payment_method.COD = true
	order_details.Price = product_details.Price

	// If the user already has products in their cart, include them in the order
	order_details.Order_cart = append(user_details.UserCart, product_details)

	// Calculate the total price for the order (sum of all products in the cart)
	totalPrice := 0
	for _, item := range order_details.Order_cart {
		totalPrice += item.Price
	}
	order_details.Price = totalPrice

	// Push the new order into the user's orders
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "orders", Value: order_details}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user collection", err)
		return err
	}

	// empty_cart := make([]models.ProductUser, 0)

	filter1 := bson.D{{Key: "_id", Value: id}}

	update1 := bson.D{{Key: "$set", Value: bson.D{
		{Key: "usercart", Value: bson.A{}},
	}}}

	_, err = userCollection.UpdateOne(ctx, filter1, update1)

	if err != nil {
		log.Println("Error emptying the cart after placing the instant order", err)
		return err
	}
	return nil
}
