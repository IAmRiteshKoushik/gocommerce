package main

import (
	controller "github.com/IAmRiteshKoushik/gocommerce/controllers"
	db "github.com/IAmRiteshKoushik/gocommerce/database"
	middleware "github.com/IAmRiteshKoushik/gocommerce/middleware"
	routes "github.com/IAmRiteshKoushik/gocommerce/routes"

	"github.com/gin-gonic/gin"

	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controller.NewApplication(db.ProductData(db.Client, "Products"), db.UserData(db.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	// Other routes
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
