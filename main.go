package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/usama1031/golang-ecommerce-store/controllers"
	"github.com/usama1031/golang-ecommerce-store/database"
	"github.com/usama1031/golang-ecommerce-store/middleware"
	"github.com/usama1031/golang-ecommerce-store/routes"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "products"), database.UserData(database.Client, "users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	router.Use(middleware.Authenicate())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/viewcart", app.ViewCart())

	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditAddress())
	// router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())

	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))

}
