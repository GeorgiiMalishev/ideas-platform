package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) (uuid.UUID, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, categoryID, coffeeShopID uuid.UUID) error
	GetByID(ctx context.Context, categoryID, coffeeShopID uuid.UUID) (models.Category, error)
	GetByCoffeeShop(ctx context.Context, coffeeShopID uuid.UUID, page, pageSize int) ([]models.Category, int, error)
}
