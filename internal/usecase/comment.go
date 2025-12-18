package usecase

import (
	"context"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/google/uuid"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, actorID, ideaID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetCommentsByIdeaID(ctx context.Context, actorID, ideaID uuid.UUID, params dto.GetCommentsRequest) ([]dto.CommentResponse, error)
	DeleteComment(ctx context.Context, actorID, ideaID, commentID uuid.UUID) error
}
