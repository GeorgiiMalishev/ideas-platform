package handlers

import (
	"log/slog"
	"net/http"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LikeHandler struct {
	usecase usecase.LikeUsecase
	logger  *slog.Logger
}

func NewLikeHandler(u usecase.LikeUsecase, l *slog.Logger) *LikeHandler {
	return &LikeHandler{
		usecase: u,
		logger:  l,
	}
}

// @Summary Like an idea
// @Description Like an idea by its ID
// @Tags likes
// @Produce json
// @Param id path string true "Idea ID"
// @Success 201 "Created"
// @Failure 400 {object} dto.ErrorResponse "Bad Request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /ideas/{id}/like [post]
// @Security ApiKeyAuth
func (h *LikeHandler) LikeIdea(c *gin.Context) {
	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid idea id"})
		return
	}

	err = h.usecase.LikeIdea(c.Request.Context(), userID, ideaID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Unlike an idea
// @Description Unlike an idea by its ID
// @Tags likes
// @Produce json
// @Param id path string true "Idea ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse "Bad Request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /ideas/{id}/unlike [delete]
// @Security ApiKeyAuth
func (h *LikeHandler) UnlikeIdea(c *gin.Context) {
	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid idea id"})
		return
	}

	err = h.usecase.UnlikeIdea(c.Request.Context(), userID, ideaID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Check if user liked an idea
// @Description Check if the current user has liked an idea by its ID
// @Tags likes
// @Produce json
// @Param id path string true "Idea ID"
// @Success 200 {object} dto.HasLikedResponse
// @Failure 400 {object} dto.ErrorResponse "Bad Request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal Server Error"
// @Router /ideas/{id}/liked [get]
// @Security ApiKeyAuth
func (h *LikeHandler) HasUserLiked(c *gin.Context) {
	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	ideaID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid idea id"})
		return
	}

	hasLiked, err := h.usecase.HasUserLiked(c.Request.Context(), userID, ideaID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, dto.HasLikedResponse{
		HasLiked: hasLiked,
	})
}
