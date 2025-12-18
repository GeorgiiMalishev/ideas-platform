package repository

import (
	"context"
	"errors"
	"fmt"

	apperrors "github.com/GeorgiiMalishev/ideas-platform/internal/app_errors"
	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.IdeaComment) (*models.IdeaComment, error) {
	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return comment, nil
}

func (r *commentRepository) GetByIdeaID(ctx context.Context, ideaID uuid.UUID, limit, offset int) ([]models.IdeaComment, error) {
	var comments []models.IdeaComment
	query := r.db.WithContext(ctx).
		Preload("Creator"). // Preload the Creator user data
		Where("idea_id = ? AND is_deleted = ?", ideaID, false).
		Order("created_at DESC") // Order by creation time, newest first

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset >= 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments by idea ID: %w", err)
	}
	return comments, nil
}

func (r *commentRepository) GetByID(ctx context.Context, commentID uuid.UUID) (*models.IdeaComment, error) {
	var comment models.IdeaComment
	if err := r.db.WithContext(ctx).Preload("Creator").First(&comment, "id = ? AND is_deleted = ?", commentID, false).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewErrNotFound("comment", commentID.String())
		}
		return nil, fmt.Errorf("failed to get comment by ID: %w", err)
	}
	return &comment, nil
}

func (r *commentRepository) Delete(ctx context.Context, commentID uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.IdeaComment{}).Where("id = ?", commentID).Update("is_deleted", true)
	if result.Error != nil {
		return fmt.Errorf("failed to delete comment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return apperrors.NewErrNotFound("comment", commentID.String())
	}
	return nil
}

func (r *commentRepository) CountByIdeaID(ctx context.Context, ideaID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.IdeaComment{}).Where("idea_id = ? AND is_deleted = ?", ideaID, false).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count comments by idea ID: %w", err)
	}
	return count, nil
}
