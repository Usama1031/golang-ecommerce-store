package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/usama1031/golang-ecommerce-store/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/add-product", controllers.AdminAddProduct())
	incomingRoutes.GET("/users/product-view", controllers.ProductViewer())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
}
