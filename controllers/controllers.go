package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"

    db "github.com/IAmRiteshKoushik/gocommerce/database"
    models "github.com/IAmRiteshKoushik/gocommerce/models"
)

var UserCollection *mongo.Collection = db.UserData(db.Client, "Users")
var ProductCollection *mongo.Collection = db.ProductData(db.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  if err != nil {
    log.Panic(err)
  }
  return string(bytes)
}

func VerifyPassword(userPassword, givenPassword string) (bool, string) {
  valid := true
  msg := ""
  err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
  if err != nil {
    msg = "Login of Password is incorrect"
    valid = false
  }

  return valid, msg 
}

func SignUp() gin.HandlerFunc {
  return func(c *gin.Context) {
    var ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    // Binding JSON with User struct
    var user models.User
    if err := c.BindJSON(&user); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{
        "error" : err.Error(),
      })
      return
    }

    // Validate user struct 
    validationErr := Validate.Struct(user)
    if validationErr != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error" : validationErr})
      return
    }

    // Check for user availability
    count, err := UserCollection.CountDocuments(ctx, bson.M{"email" : user.Email})
    if err != nil {
      log.Panic(err)
      c.JSON(http.StatusInternalServerError, gin.H{"error" : err})
      return
    }

    if count > 0 {
      c.JSON(http.StatusBadRequest, gin.H{"error" : "user already exists"})
    }

    count, err := UserCollection.CountDocuments(ctx, bson.M{"phone" : user.Phone})
    
    defer cancel()
    if err != nil {
      log.Panic(err)
      c.JSON(http.StatusInternalServerError, gin.H{"error" : err})
      return
    }

    if count > 0 {
      c.JSON(http.StatusBadRequest, gin.H{"error" : "this phone no. already in use"})
      return
    }

    password := HashPassword(*user.Password)
    user.Password = &password

    user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
    user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
    user.ID = primitive.NewObjectID()
    user.User_ID = user.ID.Hex()

token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_name, *user.Last_name, user.User_ID)
    user.Token = &token
    user.Refresh_Token = &refreshtoken
    user.UserCart = make([]models.ProductUser, 0)
    user.Address_Details = make([]models.Address, 0)
    user.Order_Status = make([]models.Order, 0)
    
    _, insertErr := UserCollection.InsertOne(ctx, user)
    if insertErr != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error":"the user did not get created"})
      return
    }
    defer cancel()

    // Creationg confirmation message
    c.JSON(http.StatusCreated, "Successfully signed in!")

  }
}

func Login() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

    var user models.User
    if err := c.BindJSON(&user); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error" : err})
      return
    }

    // Check if the user even exists in the database
    err := UserCollection.FindOne(ctx, bson.M{"email" : user.Email}).Decode(&foundUser)
    defer cancel()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error" : "Login or password incorrect"})
      return
    }

    PasswordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
    defer cancel()
    if !PasswordIsValid {
      c.JSON(http.StatusInternalServerError, gin.H{"error" : msg})
      fmt.Println(msg)
      return
    }

    token, refreshToken, _ := generate.TOkenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_name, *foundUser.User_ID)
    defer cancel()

    generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
    c.JSON(http.StatusFound, foundUser)
  }
}

func ProductViewerAdmin() gin.HandlerFunc {
  return func(c *gin.Context){
    
  }
}

func SearchProduct() gin.HandlerFunc {
  return func(c *gin.Context){
    var productList []models.Product 
    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

    cursor, err := ProductCollection.Find(ctx, bson.D{{}})
    if err != nil {
      c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, try after some time.")
      return
    }

    err = cursor.All(ctx, &productList)
    if err != nil{
      log.Println(err)
      c.AbortWithStatus(http.StatusInternalServerError)
      return
    }

    defer cursor.Close()
    if err := cursor.Err(); err != nil {
      log.Println(err)
      c.IndentedJSON(http.StatusBadRequest, "invalid")
      return
    }
    defer cancel()
    c.IndentedJSON(http.StatusOK, productList)
  }
}

func SearchProductByQuery() gin.HandlerFunc {
  return func(c *gin.Context){
    var searchProduct []models.Product
    queryParam := c.Query("name")

    // You want to check if it's empty
    if queryParam == "" {
      log.Println("Query is empty")
      c.Header("Context-type", "application/json")
      c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
      c.Abort()
      return
    }

    var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
    defer cancel()

    searchQueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex":queryParam}})
    if err != nil {
      c.IndentedJSON(http.StatusNotFound, "Something went wrong while fetching the data")
      return
    }
    err = searchQueryDB.All(ctx, &searchProduct)
    if err != nil {
      log.Println(err)
      c.IndentedJSON(http.StatusBadRequest, "Invalid")
      return
    }

    defer searchQueryDB.Close()

    if err := searchQueryDB.Err(); err != nil {
      log.Println(err)
      c.IndentedJSON(http.StatusBadRequest, "Invalid request")
      return
    }
    defer cancel()

    c.IndentedJSON(http.StatusOK, searchProduct)
  }

}
