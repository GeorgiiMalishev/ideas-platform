package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	uc     usecase.CommentUsecase
	logger *slog.Logger
}

func NewCommentHandler(uc usecase.CommentUsecase, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{
		uc:     uc,
		logger: logger,
	}
}

// @Summary Create a new comment for an idea
// @Description Create a new comment for a specific idea
// @Tags ideas
// @Accept json
// @Produce json
// @Param id path string true "Idea ID"
// @Param comment body dto.CreateCommentRequest true "Comment information"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ideas/{id}/comments [post]
// @Security ApiKeyAuth
func (h *CommentHandler) CreateComment(c *gin.Context) {
	ideaID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind create comment request", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	resp, err := h.uc.CreateComment(c.Request.Context(), actorID, ideaID, &req)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Get comments for an idea
// @Description Get a list of comments for a specific idea with optional pagination
// @Tags ideas
// @Produce json
// @Param id path string true "Idea ID"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.CommentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ideas/{id}/comments [get]
// @Security ApiKeyAuth
func (h *CommentHandler) GetComments(c *gin.Context) {
	ideaID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}

	pageRaw := c.Query("page")
	limitRaw := c.Query("limit")
	// Sort is not used in the current usecase implementation but can be added later
	// sort := c.Query("sort")

	page, _ := strconv.Atoi(pageRaw)
	limit, _ := strconv.Atoi(limitRaw)

	params := dto.GetCommentsRequest{
		Page:  page,
		Limit: limit,
		// Sort:  sort,
	}

	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	resp, err := h.uc.GetCommentsByIdeaID(c.Request.Context(), actorID, ideaID, params)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Delete a comment
// @Description Delete a specific comment by ID for a given idea
// @Tags ideas
// @Produce json
// @Param id path string true "Idea ID"
// @Param comment_id path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ideas/{id}/comments/{comment_id} [delete]
// @Security ApiKeyAuth
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	ideaID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}

	commentIDStr := c.Param("comment_id")
	commentID, err := uuid.Parse(commentIDStr)
	if err != nil {
		h.logger.Error("failed to parse commentID", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	err = h.uc.DeleteComment(c.Request.Context(), actorID, ideaID, commentID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}
