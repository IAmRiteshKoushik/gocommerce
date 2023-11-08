package database

import "errors"

// Error variables
var (
  
  ErrCantFindProduct = errors.New("Cannot find the product.")
  ErrCantDecodeProduct = errors.New("Cannot find the product.") 
  ErrUserIdIsNotValid = errors.New("This user is not valid.")
  ErrCantUpdateUser = errors.New("Cannot add this product to the cart.")
  ErrCantRemoveItemCart = errors.New("Cannot remove this item from the cart.")
  ErrCantGetItem = errors.New("Was unable to get the item from the cart.")
  ErrCantBuyCartItem = errors.New("Cannot update the purchase.")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
