package repository

import (
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
)

type WorkerCoffeeShopRepository interface {
	Create(workerShop *models.WorkerCoffeeShop) (*models.WorkerCoffeeShop, error)
	GetByID(id uuid.UUID) (*models.WorkerCoffeeShop, error)
	ListByCoffeeShopID(coffeeShopID uuid.UUID, limit, offset int) ([]models.WorkerCoffeeShop, error)
	ListByWorkerID(workerID uuid.UUID, limit, offset int) ([]models.WorkerCoffeeShop, error)
	Update(workerShop *models.WorkerCoffeeShop) error
	Delete(id uuid.UUID) error
	IsWorkerInShop(workerID, shopID uuid.UUID) (bool, error)
}
