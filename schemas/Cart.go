package schemas

import "go-server/models"

type AddToCartRequest struct {
    ProductCardID uint `json:"product_card_id" validate:"required"`
    Quantity     int  `json:"quantity" validate:"required,gt=0"`
}

type UpdateCartItemRequest struct {
    Quantity int `json:"quantity" validate:"required,gt=0"`
}

type CartResponse struct {
    ID       uint               `json:"id"`
    UserID   uint               `json:"user_id"`
    Items    []CartItemResponse `json:"items"`
    Total    int                `json:"total"`
}

type CartItemResponse struct {
    ID           uint              `json:"id"`
    ProductCard  models.ProductCard `json:"product"`
    Quantity     int               `json:"quantity"`
    Subtotal     int               `json:"subtotal"`
}