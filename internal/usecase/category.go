package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
)

type CategoryUsecase interface {
	Create(ctx context.Context, userID uuid.UUID, coffeeShopID uuid.UUID, category dto.CreateCategory) (uuid.UUID, error)
	Update(ctx context.Context, userID, coffeeShopID, categoryID uuid.UUID, category dto.UpdateCategory) error
	Delete(ctx context.Context, userID uuid.UUID, coffeeShopID uuid.UUID, categoryID uuid.UUID) error
	GetByID(ctx context.Context, coffeeShopID, categoryID uuid.UUID) (dto.CategoryResponse, error)
	GetByCoffeeShop(ctx context.Context, coffeeShopID uuid.UUID, page, pageSize int) ([]dto.CategoryResponse, int, error)
}
