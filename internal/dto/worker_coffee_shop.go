package dto

import "github.com/google/uuid"

// AddWorkerToShopRequest defines the request body for adding a worker to a coffee shop.
type AddWorkerToShopRequest struct {
	WorkerID     uuid.UUID `json:"worker_id" binding:"required"`
	CoffeeShopID uuid.UUID `json:"coffee_shop_id" binding:"required"`
}

// WorkerCoffeeShopResponse defines the response for a worker-coffeeshop relationship.
type WorkerCoffeeShopResponse struct {
	ID         uuid.UUID          `json:"id"`
	Worker     UserResponse       `json:"worker"`
	CoffeeShop CoffeeShopResponse `json:"coffee_shop"`
}
