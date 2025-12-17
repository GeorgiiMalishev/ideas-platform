package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/GeorgiiMalishev/ideas-platform/internal/app_errors"
	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/GeorgiiMalishev/ideas-platform/internal/repository"
	"gorm.io/gorm"
)

type CategoryUsecaseImpl struct {
	categoryRepo  repository.CategoryRepository
	accessControl AccessControlUsecase
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository, accessControl AccessControlUsecase) CategoryUsecase {
	return &CategoryUsecaseImpl{categoryRepo: categoryRepo, accessControl: accessControl}
}

func (u *CategoryUsecaseImpl) Create(ctx context.Context, userID, coffeeShopID uuid.UUID, category dto.CreateCategory) (uuid.UUID, error) {
	if err := u.accessControl.CanManageCoffeeShop(ctx, userID, coffeeShopID); err != nil {
		return uuid.Nil, err
	}

	newCategory := &models.Category{
		CoffeeShopID: &coffeeShopID,
		Title:        category.Title,
		Description:  category.Description,
	}

	return u.categoryRepo.Create(ctx, newCategory)
}

func (u *CategoryUsecaseImpl) Update(ctx context.Context, userID, coffeeShopID, categoryID uuid.UUID, category dto.UpdateCategory) error {
	if err := u.accessControl.CanManageCoffeeShop(ctx, userID, coffeeShopID); err != nil {
		return err
	}

	updateCategory := &models.Category{
		ID:           categoryID,
		CoffeeShopID: &coffeeShopID,
		Title:        category.Title,
		Description:  category.Description,
	}

	return u.categoryRepo.Update(ctx, updateCategory)
}

func (u *CategoryUsecaseImpl) Delete(ctx context.Context, userID, coffeeShopID, categoryID uuid.UUID) error {
	if err := u.accessControl.CanManageCoffeeShop(ctx, userID, coffeeShopID); err != nil {
		return err
	}

	return u.categoryRepo.Delete(ctx, categoryID, coffeeShopID)
}

func (u *CategoryUsecaseImpl) GetByID(ctx context.Context, coffeeShopID, categoryID uuid.UUID) (dto.CategoryResponse, error) {
	category, err := u.categoryRepo.GetByID(ctx, categoryID, coffeeShopID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.CategoryResponse{}, apperrors.NewErrNotFound("category", categoryID.String())
		}
		return dto.CategoryResponse{}, fmt.Errorf("failed to get category by id: %w", err)
	}

	return dto.CategoryResponse{
		ID:           category.ID,
		CoffeeShopID: category.CoffeeShopID,
		Title:        category.Title,
		Description:  category.Description,
	}, nil
}

func (u *CategoryUsecaseImpl) GetByCoffeeShop(ctx context.Context, coffeeShopID uuid.UUID, page, pageSize int) ([]dto.CategoryResponse, int, error) {
	categories, total, err := u.categoryRepo.GetByCoffeeShop(ctx, coffeeShopID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get categories by coffee shop: %w", err)
	}

	var categoryResponses []dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, dto.CategoryResponse{
			ID:           category.ID,
			CoffeeShopID: category.CoffeeShopID,
			Title:        category.Title,
			Description:  category.Description,
		})
	}

	return categoryResponses, total, nil
}
