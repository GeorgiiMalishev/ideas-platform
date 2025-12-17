package usecase

import (
	"context"
	"log/slog"

	"github.com/GeorgiiMalishev/ideas-platform/internal/repository"
	"github.com/google/uuid"
)

type AccessControlUsecaseImpl struct {
	workerShopRepo repository.WorkerCoffeeShopRepository
	logger         *slog.Logger
}

func NewAccessControlUsecase(workerShopRepo repository.WorkerCoffeeShopRepository, logger *slog.Logger) AccessControlUsecase {
	return &AccessControlUsecaseImpl{
		workerShopRepo: workerShopRepo,
		logger:         logger,
	}
}

func (u *AccessControlUsecaseImpl) CanManageCoffeeShop(ctx context.Context, userID, coffeeShopID uuid.UUID) error {
	return CheckShopAdminAccess(ctx, u.logger, u.workerShopRepo, userID, coffeeShopID)
}
