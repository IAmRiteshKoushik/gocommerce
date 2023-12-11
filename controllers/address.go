package controllers

import (
    "fmt"
	"net/http"
	"time"
    "context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

	models "github.com/IAmRiteshKoushik/gocommerce/models"
)


func AddAddress() gin.HandlerFunc {
  return func(c *gin.Context){
    user_id := c.Query("id")

    if user_id == ""{
      c.Header("Content-type", "application/json")
      c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
      c.Abort()
      return
    }
    
    address, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
    }

    var addresses models.Address

    addresses.Address_id = primitive.NewObjectID()
    if err := c.BindJSON(&addresses); err != nil {
      c.IndentedJSON(http.StatusNotAcceptable, err.Error())
    }

    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)

    // Aggregation stages
    match_filter := bson.D{
      Key: "$match",
      Value: bson.D{primitive.E{
        Key: "_id",
        Value: address,
      }},
    }
    unwind := bson.D{
      Key: "$unwind",
      Value: bson.D{primitive.E{
        Key: "path",
        Value: "$address",
      }},
    }
    grouping := bson.D{
      Key: "$group",
      Value: bson.D{primitive.E{
        Key: "_id",
        Value: "$address_id",
      }},
      {
        Key: "count",
        Value: bson.D{primitive.E{
          Key: "$sum",
          Value: 1,
        }},
      },
    }

    pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, grouping})
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
    }

    var addressInfo []bson.M
    if err := pointCursor.All(ctx, &addressInfo); err != nil {
      panic(err)
    }

    var size int32
    for _, addressNo := range addressInfo {
      count := addressNo["count"]
      size = count.(int32)
    }
    if size < 2 {
      filter := bson.D{primitive.E{
        Key: "_id",
        Value: address,
      }}
      update := bson.D{{
        Key: "$push",
        Value: bson.D{primitive.E{
          Key: "address",
          Value: addresses,
        }},
      }}
      _, err := UserCollection.UpdateOne(ctx, filter, update)
      if err != nil {
        fmt.Println(err)
        }
      } else {
      c.IndentedJSON(http.StatusBadRequest, "Not Allowed")
    }
    defer cancel()
    ctx.Done()
  }
}

func EditWorkAddress() gin.HandlerFunc {
  return func(c *gin.Context){
    user_id := c.Query("id")

    if user_id == ""{
      c.Header("Content-type", "application/json")
      c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
      c.Abort()
      return
    }

    usert_id, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
    }

    var editaddress models.Address 
    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()
    filter := bson.D{primitive.E{Key:"_id", Value:usert_id}}
    update := bson.D{{
      Key: "$set",
      Value: bson.D{primitive.E{
        Key:"address.1.house_name",
        Value: editaddress.House,
      },
      {
        Key: "address.1.street_name",
        Value: editaddress.Street,
      },
      {
        Key: "address.1.city_name",
        Value: editaddress.City,
      },
      {
        Key: "address.1.pin_code",
        Value: editaddress.Pincode,
      },
    }}}

    _, err = UserCollection.UpdateByID(ctx, filter, update)
    if err != nil {
      c.IndentedJSON(http.StatusBadRequest, "Something went wrong")
    }
    defer cancel()
    ctx.Done()
    c.IndentedJSON(http.StatusOK, "Successfully the updated work address")


  }
}

func EditHouseAddress() gin.HandlerFunc {
  return func(c *gin.Context){
    user_id := c.Query("id")

    if user_id == ""{
      c.Header("Content-type", "application/json")
      c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
      c.Abort()
      return
    }

    var editaddress models.Address 
    if err := c.BindJSON(&editaddress); err != nil {
      c.IndentedJSON(http.StatusBadRequest, err.Error())
    }
    usert_id, err := primitive.ObjectIDFromHex(user_id)
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Internal Server Error")
    }

    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()
    filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}
    update := bson.D{{
      Key: "$set",
      Value: bson.D{primitive.E{
        Key:"address.0.house_name",
        Value: editaddress.House,
      },
      {
        Key: "address.0.street_name",
        Value: editaddress.Street,
      },
      {
        Key: "address.0.city_name",
        Value: editaddress.City,
      },
      {
        Key: "address.0.pin_code",
        Value: editaddress.Pincode,
      },
    }}}

    _, err = UserCollection.UpdateOne(ctx, filter, update)
    if err != nil {
      c.IndentedJSON(http.StatusBadRequest, "Something went wrong")
      return
    }
    defer cancel()
    ctx.Done()
    c.IndentedJSON(http.StatusOK, "Successfully update the home address")
  }
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
