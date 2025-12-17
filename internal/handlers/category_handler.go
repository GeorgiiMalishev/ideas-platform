package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
	logger          *slog.Logger
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase, logger: logger}
}

// @Summary Create a new category
// @Description Create a new category for a coffee shop
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param category body dto.CreateCategory true "Category information"
// @Success 201 {object} map[string]string
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/categories [post]
// @Security ApiKeyAuth
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind category create request: ", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	coffeeShopID, ok := parseUUIDFromParam(h.logger, c, "id")
	if !ok {
		return
	}

	id, err := h.categoryUsecase.Create(c.Request.Context(), userID, coffeeShopID, req)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Update a category
// @Description Update category details for the given ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param category_id path string true "Category ID"
// @Param category body dto.UpdateCategory true "Category update information"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/categories/{category_id} [put]
// @Security ApiKeyAuth
func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.UpdateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind category update request: ", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	coffeeShopID, ok := parseUUIDFromParam(h.logger, c, "id")
	if !ok {
		return
	}

	categoryID, ok := parseUUIDFromParam(h.logger, c, "category_id")
	if !ok {
		return
	}

	if err := h.categoryUsecase.Update(c.Request.Context(), userID, coffeeShopID, categoryID, req); err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param category_id path string true "Category ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/categories/{category_id} [delete]
// @Security ApiKeyAuth
func (h *CategoryHandler) Delete(c *gin.Context) {
	userID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}

	coffeeShopID, ok := parseUUIDFromParam(h.logger, c, "id")
	if !ok {
		return
	}

	categoryID, ok := parseUUIDFromParam(h.logger, c, "category_id")
	if !ok {
		return
	}

	if err := h.categoryUsecase.Delete(c.Request.Context(), userID, coffeeShopID, categoryID); err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param category_id path string true "Category ID"
// @Success 200 {object} dto.CategoryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/categories/{category_id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	coffeeShopID, ok := parseUUIDFromParam(h.logger, c, "id")
	if !ok {
		return
	}

	categoryID, ok := parseUUIDFromParam(h.logger, c, "category_id")
	if !ok {
		return
	}

	category, err := h.categoryUsecase.GetByID(c.Request.Context(), coffeeShopID, categoryID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Get categories by coffee shop
// @Description Get a list of all categories for a given coffee shop with optional pagination
// @Tags categories
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/categories [get]
func (h *CategoryHandler) GetByCoffeeShop(c *gin.Context) {
	coffeeShopID, ok := parseUUIDFromParam(h.logger, c, "id")
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	categories, total, err := h.categoryUsecase.GetByCoffeeShop(c.Request.Context(), coffeeShopID, page, pageSize)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories, "total": total})
}
