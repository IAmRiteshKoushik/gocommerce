package routes

import (
	controller "github.com/IAmRiteshKoushik/gocommerce/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.SignUp())
	incomingRoutes.POST("/users/login", controller.Login())
	incomingRoutes.POST("/admin/addproduct", controller.ProductViewerAdmin())
	incomingRoutes.GET("/users/search", controller.SearchProduct())
	incomingRoutes.GET("/users/search", controller.SearchProductByQuery())
}
