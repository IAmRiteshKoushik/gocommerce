package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/IAmRiteshKoushik/gocommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Error variables
var (
	ErrCantFindProduct    = errors.New("Cannot find the product.")
	ErrCantDecodeProduct  = errors.New("Cannot find the product.")
	ErrUserIdIsNotValid   = errors.New("This user is not valid.")
	ErrCantUpdateUser     = errors.New("Cannot add this product to the cart.")
	ErrCantRemoveItemCart = errors.New("Cannot remove this item from the cart.")
	ErrCantGetItem        = errors.New("Was unable to get the item from the cart.")
	ErrCantBuyCartItem    = errors.New("Cannot update the purchase.")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchfromdb, err := prodCollection.Find(ctx, bson.M{
		"_id": productID,
	})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var productCart []models.ProductUser
	err = searchfromdb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{
		Key:   "_id",
		Value: id,
	}}
	update := bson.D{primitive.E{
		Key: "$push",
		Value: bson.D{
			primitive.E{
				Key:   "usercart",
				Value: bson.D{{Key: "$each", Value: productCart}},
			},
		},
	}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{
		Key:   "_id",
		Value: id,
	}}
	update := bson.M{
		"$pull": bson.M{
			"usercart": bson.M{
				"_id": productID,
			},
		},
	}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	// fetch the cart of the user
	// Find the cart total
	// Create an order with the items
	// added order to the user colleciton
	// added items in the cart to user collection
	// Empty up the cart

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Ordered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{{
		Key: "$unwind",
		Value: bson.D{primitive.E{
			Key:   "path",
			Value: "$usercart",
		}},
	}}
	grouping := bson.D{{
		Key: "$group",
		Value: bson.D{primitive.E{
			Key:   "path",
			Value: "$_id",
		},
			{
				Key: "total",
				Value: bson.D{primitive.E{
					Key:   "$sum",
					Value: "$usercart.price",
				}},
			},
		},
	}}
	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		panic(err)
	}

	var getUserCart []bson.M
	if err := currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
	}

	var total_price int32
	for _, user_item := range getUserCart {
		price := user_item["total"]
		total_price = price.(int32)
	}

	orderCart.Price = int(total_price)
	filter := bson.D{primitive.E{
		Key:   "_id",
		Value: id,
	}}
	update := bson.D{{
		Key: "$push",
		Value: bson.D{primitive.E{
			Key:   "orders",
			Value: orderCart,
		}},
	}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{
		Key:   "_id",
		Value: id,
	}}).Decode(&getCartItems)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{
		"$push": bson.M{
			"orders.$[].order_list": bson.M{
				"$each": getCartItems.UserCart,
			},
		},
	}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	usercart_empty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{
		Key: "$set",
		Value: bson.D{primitive.E{
			Key:   "usercart",
			Value: usercart_empty,
		}},
	}}}}
	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, userCollection, prodCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var product_detail models.ProductUser
	var orders_detail models.Order

	orders_detail.Order_ID = primitive.NewObjectID()
	orders_detail.Ordered_At = time.Now()
	orders_detail.Order_Cart = make([]models.ProductUser, 0)
	orders_detail.Payment_Method.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{
		Key:   "_id",
		Value: productID,
	}}).Decode(&product_detail)
	if err != nil {
		log.Println(err)
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{
		"$push": bson.M{
			"orders.$[].order_list": product_detail,
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}

	return nil
}
