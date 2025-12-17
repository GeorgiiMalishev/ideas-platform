package dto

import "github.com/google/uuid"

type CreateCategory struct {
	Title       string  `json:"title" binding:"required,min=3,max=50"`
	Description *string `json:"description"`
}

type UpdateCategory struct {
	Title       string  `json:"title" binding:"required,min=3,max=50"`
	Description *string `json:"description"`
}

type CategoryResponse struct {
	ID           uuid.UUID  `json:"id"`
	CoffeeShopID *uuid.UUID `json:"coffee_shop_id"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
}
