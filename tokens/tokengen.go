package tokens

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	db "github.com/IAmRiteshKoushik/gocommerce/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string `json:"email"`
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var UserData *mongo.Collection = db.UserData(db.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email, firstname, lastname, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is already expired"
		return
	}

	return claims, msg

}

func UpdateAllTokens(signedToken, signedRefreshToken, userid string) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{
		Key:   "token",
		Value: signedToken,
	})
	updateObj = append(updateObj, bson.E{
		Key:   "refresh_token",
		Value: signedRefreshToken,
	})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{
		Key:   "updateat",
		Value: updatedAt,
	})

	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
	}

}
