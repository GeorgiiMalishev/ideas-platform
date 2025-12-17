package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return uuid.Nil, fmt.Errorf("failed to create category: %w", err)
	}
	return category.ID, nil
}

func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ? AND coffee_shop_id = ?", category.ID, category.CoffeeShopID).Updates(category).Error
}

func (r *CategoryRepositoryImpl) Delete(ctx context.Context, categoryID, coffeeShopID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ? AND coffee_shop_id = ?", categoryID, coffeeShopID).Update("is_deleted", true).Error
}

func (r *CategoryRepositoryImpl) GetByID(ctx context.Context, categoryID, coffeeShopID uuid.UUID) (models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, "id = ? AND coffee_shop_id = ? AND is_deleted = false", categoryID, coffeeShopID).Error
	return category, err
}

func (r *CategoryRepositoryImpl) GetByCoffeeShop(ctx context.Context, coffeeShopID uuid.UUID, page, pageSize int) ([]models.Category, int, error) {
	var categories []models.Category
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.WithContext(ctx).Model(&models.Category{}).
		Where("coffee_shop_id = ? AND is_deleted = false", coffeeShopID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Where("coffee_shop_id = ? AND is_deleted = false", coffeeShopID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, int(total), nil
}
