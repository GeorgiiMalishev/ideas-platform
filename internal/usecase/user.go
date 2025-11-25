package usecase

import (
	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/google/uuid"
)

type UserUsecase interface {
	UpdateUser(actorID, ID uuid.UUID, req *dto.UpdateUserRequest) error
	GetAllUsers(actorID uuid.UUID, page, limit int) ([]dto.UserResponse, error)
	GetUser(actorID, ID uuid.UUID) (*dto.UserResponse, error)
	DeleteUser(actorID, ID uuid.UUID) error
}
