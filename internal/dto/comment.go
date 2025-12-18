package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	Text       string `json:"text" binding:"required"`
	AuthorName string `json:"name" binding:"required"`
}

type CommentResponse struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	AuthorName string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetCommentsRequest struct {
	Page  int
	Limit int
	Sort  string
}
