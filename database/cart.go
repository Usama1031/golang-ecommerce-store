package database

import "errors"

var (
	ErrCannotFindProduct        = errors.New("cannot find the product")
	ErrCannotDecodeProducts     = errors.New("cannot find the product")
	ErrUserIdIsNotValid         = errors.New("this user is not valid")
	ErrCannotUpdateUser         = errors.New("cannot add this product to the cart")
	ErrCannotRemoveItemFromCart = errors.New("cannot remove this item from the cart")
	ErrCannotGetItem            = errors.New("unable to get the item from cart")
	ErrCannotBuyCartItem        = errors.New("cannot update the purchase")
)

func AddProductToCart() {

}

func RemoveItemFromCart() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
