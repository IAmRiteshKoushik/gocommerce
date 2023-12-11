package controllers

import (
	"net/http"
	"time"
    "context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	models "github.com/IAmRiteshKoushik/gocommerce/models"
)

func AddAddress() gin.HandlerFunc {

}

func EditHomeAddress() gin.HandlerFunc {

}

func EditWorkAddress() gin.HandlerFunc {

}

func DeleteAddress() gin.HandlerFunc {
  return func(c *gin.Context){
    user_id := c.Query("id")

    if user_id == ""{
      c.Header("Content-type", "application/json")
      c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
      c.Abort()
      return
    }

    address := make([]models.Address, 0)
    usert_id, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
    }

    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()
    filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}
    update := bson.D{{Key:"$set", Value: bson.D{primitive.E{Key:"address", Value: address}}}}

    _, err = UserCollection.UpdateOne(ctx, filter, update)
    if err != nil {
      c.IndentedJSON(http.StatusNotFound, "Wrong command")
      return
    }
    defer cancel()
    c.IndentedJSON(http.StatusOK, "Successfully deleted")
  }
}
