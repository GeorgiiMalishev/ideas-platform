package repository

import (
	"context"

	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.IdeaComment) (*models.IdeaComment, error)
	GetByIdeaID(ctx context.Context, ideaID uuid.UUID, limit, offset int) ([]models.IdeaComment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*models.IdeaComment, error)
	Delete(ctx context.Context, commentID uuid.UUID) error
	CountByIdeaID(ctx context.Context, ideaID uuid.UUID) (int64, error)
}
