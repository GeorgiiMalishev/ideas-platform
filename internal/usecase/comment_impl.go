package usecase

import (
	"context"
	"errors"
	"log/slog"

	apperrors "github.com/GeorgiiMalishev/ideas-platform/internal/app_errors"
	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/GeorgiiMalishev/ideas-platform/internal/repository"
	"github.com/google/uuid"
)

type commentUsecase struct {
	commentRepo          repository.CommentRepository
	ideaRepo             repository.IdeaRepository
	workerCoffeeShopRepo repository.WorkerCoffeeShopRepository
	logger               *slog.Logger
}

func NewCommentUsecase(
	commentRepo repository.CommentRepository,
	ideaRepo repository.IdeaRepository,
	workerCoffeeShopRepo repository.WorkerCoffeeShopRepository,
	logger *slog.Logger,
) CommentUsecase {
	return &commentUsecase{
		commentRepo:          commentRepo,
		ideaRepo:             ideaRepo,
		workerCoffeeShopRepo: workerCoffeeShopRepo,
		logger:               logger,
	}
}

func (uc *commentUsecase) CreateComment(ctx context.Context, actorID, ideaID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	l := uc.logger.With("method", "CreateComment", "actorID", actorID, "ideaID", ideaID)

	idea, err := uc.ideaRepo.GetIdea(ctx, ideaID)
	if err != nil {
		l.Error("failed to get idea", slog.String("error", err.Error()))
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			return nil, apperrors.NewErrNotFound("idea", ideaID.String())
		}
		return nil, err
	}
	if idea.CoffeeShopID == nil {
		l.Error("idea has no associated coffee shop ID", slog.Any("idea", idea))
		return nil, errors.New("idea is not associated with a coffee shop")
	}

	// Check if the actor is a worker in the coffee shop associated with the idea
	_, err = uc.workerCoffeeShopRepo.GetByUserIDAndShopID(ctx, actorID, *idea.CoffeeShopID)
	if err != nil {
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			l.Warn("access denied: user is not worker for this coffee shop", slog.String("error", err.Error()))
			return nil, apperrors.NewErrAccessDenied("user is not a worker for this coffee shop")
		}
		l.Error("failed to check worker status", slog.String("error", err.Error()))
		return nil, err
	}

	comment := &models.IdeaComment{
		CreatorID:  &actorID,
		IdeaID:     &ideaID,
		Text:       req.Text,
		AuthorName: req.AuthorName,
	}

	createdComment, err := uc.commentRepo.Create(ctx, comment)
	if err != nil {
		l.Error("failed to create comment", slog.String("error", err.Error()))
		return nil, err
	}

	return &dto.CommentResponse{
		ID:         createdComment.ID,
		Text:       createdComment.Text,
		AuthorName: createdComment.AuthorName,
		CreatedAt:  createdComment.CreatedAt,
	}, nil
}

func (uc *commentUsecase) GetCommentsByIdeaID(ctx context.Context, actorID, ideaID uuid.UUID, params dto.GetCommentsRequest) ([]dto.CommentResponse, error) {
	l := uc.logger.With("method", "GetCommentsByIdeaID", "actorID", actorID, "ideaID", ideaID)

	idea, err := uc.ideaRepo.GetIdea(ctx, ideaID)
	if err != nil {
		l.Error("failed to get idea", slog.String("error", err.Error()))
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			return nil, apperrors.NewErrNotFound("idea", ideaID.String())
		}
		return nil, err
	}
	if idea.CoffeeShopID == nil {
		l.Error("idea has no associated coffee shop ID", slog.Any("idea", idea))
		return nil, errors.New("idea is not associated with a coffee shop")
	}

	// Check if the actor is a worker in the coffee shop associated with the idea
	_, err = uc.workerCoffeeShopRepo.GetByUserIDAndShopID(ctx, actorID, *idea.CoffeeShopID)
	if err != nil {
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			l.Warn("access denied: user is not worker for this coffee shop", slog.String("error", err.Error()))
			return nil, apperrors.NewErrAccessDenied("user is not a worker for this coffee shop")
		}
		l.Error("failed to check worker status", slog.String("error", err.Error()))
		return nil, err
	}

	limit, offset := calculatePagination(params.Page, params.Limit)

	comments, err := uc.commentRepo.GetByIdeaID(ctx, ideaID, limit, offset)
	if err != nil {
		l.Error("failed to get comments by idea ID", slog.String("error", err.Error()))
		return nil, err
	}

	var responses []dto.CommentResponse
	for _, comment := range comments {
		responses = append(responses, dto.CommentResponse{
			ID:         comment.ID,
			Text:       comment.Text,
			AuthorName: comment.AuthorName,
			CreatedAt:  comment.CreatedAt,
		})
	}
	return responses, nil
}

func (uc *commentUsecase) DeleteComment(ctx context.Context, actorID, ideaID, commentID uuid.UUID) error {
	l := uc.logger.With("method", "DeleteComment", "actorID", actorID, "ideaID", ideaID, "commentID", commentID)

	idea, err := uc.ideaRepo.GetIdea(ctx, ideaID)
	if err != nil {
		l.Error("failed to get idea", slog.String("error", err.Error()))
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			return apperrors.NewErrNotFound("idea", ideaID.String())
		}
		return err
	}
	if idea.CoffeeShopID == nil {
		l.Error("idea has no associated coffee shop ID", slog.Any("idea", idea))
		return errors.New("idea is not associated with a coffee shop")
	}

	// Check if the actor is a worker in the coffee shop associated with the idea
	_, err = uc.workerCoffeeShopRepo.GetByUserIDAndShopID(ctx, actorID, *idea.CoffeeShopID)
	if err != nil {
		var errNotFound *apperrors.ErrNotFound
		if errors.As(err, &errNotFound) {
			l.Warn("access denied: user is not worker for this coffee shop", slog.String("error", err.Error()))
			return apperrors.NewErrAccessDenied("user is not a worker for this coffee shop")
		}
		l.Error("failed to check worker status", slog.String("error", err.Error()))
		return err
	}

	// Check if the comment exists and belongs to the idea
	comment, err := uc.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		l.Error("failed to get comment by ID", slog.String("error", err.Error()))
		return err
	}
	if comment.IdeaID == nil || *comment.IdeaID != ideaID {
		l.Warn("comment does not belong to the specified idea", slog.String("commentID", commentID.String()), slog.String("ideaID", ideaID.String()))
		return apperrors.NewErrNotValid("comment does not belong to this idea")
	}

	// Any worker of the coffee shop can delete a comment
	if err := uc.commentRepo.Delete(ctx, commentID); err != nil {
		l.Error("failed to delete comment", slog.String("error", err.Error()))
		return err
	}

	return nil
}

// calculatePagination helper function
func calculatePagination(page, limit int) (int, int) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if page <= 0 {
		page = 1 // Default page
	}
	offset := (page - 1) * limit
	return limit, offset
}
