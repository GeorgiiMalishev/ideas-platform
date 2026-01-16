package dto

import "github.com/google/uuid"

type LikeRequest struct {
	IdeaID uuid.UUID `json:"idea_id"`
}

type HasLikedResponse struct {
	HasLiked bool `json:"has_liked"`
}
