package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type RewardTypeHandler struct {
	uc     usecase.RewardTypeUsecase
	logger *slog.Logger
}

func NewRewardTypeHandler(uc usecase.RewardTypeUsecase, logger *slog.Logger) *RewardTypeHandler {
	return &RewardTypeHandler{
		uc:     uc,
		logger: logger,
	}
}

// @Summary Get reward type by ID
// @Description Get reward type details by ID
// @Tags rewards
// @Produce json
// @Param id path string true "Reward Type ID"
// @Success 200 {object} dto.RewardTypeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /rewards/type/{id} [get]
// @Security ApiKeyAuth
func (h *RewardTypeHandler) GetRewardType(c *gin.Context) {
	rewardTypeID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}
	rewardType, err := h.uc.GetRewardType(c.Request.Context(), rewardTypeID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, rewardType)
}

// @Summary Get reward types by coffee shop
// @Description Get a list of all reward types for a given coffee shop with optional pagination
// @Tags rewards
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RewardTypeResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /coffee-shops/{id}/rewards/type [get]
// @Security ApiKeyAuth
func (h *RewardTypeHandler) GetRewardTypesByCoffeeShop(c *gin.Context) {
	coffeeShopID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}
	pageRaw := c.Query("page")
	limitRaw := c.Query("limit")

	page, _ := strconv.Atoi(pageRaw)
	limit, _ := strconv.Atoi(limitRaw)

	rewardTypes, err := h.uc.GetRewardsTypesFromCoffeeShop(c.Request.Context(), coffeeShopID, page, limit)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, rewardTypes)
}

// @Summary Create a new reward type
// @Description Create a new reward type
// @Tags rewards
// @Accept json
// @Produce json
// @Param reward_type body dto.CreateRewardTypeRequest true "Reward type information"
// @Success 201 {object} dto.RewardTypeResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/rewards/type [post]
// @Security ApiKeyAuth
func (h *RewardTypeHandler) CreateRewardType(c *gin.Context) {
	var req dto.CreateRewardTypeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{Message: "bad request"})
		return
	}
	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}
	resp, err := h.uc.CreateRewardType(c.Request.Context(), actorID, &req)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Update reward type by ID
// @Description Update reward type details for the given ID
// @Tags rewards
// @Accept json
// @Produce json
// @Param id path string true "Reward Type ID"
// @Param reward_type body dto.UpdateRewardTypeRequest true "Reward type update information"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/rewards/type/{id} [put]
// @Security ApiKeyAuth
func (h *RewardTypeHandler) UpdateRewardType(c *gin.Context) {
	rewardTypeID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}
	var req dto.UpdateRewardTypeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ErrorResponse{Message: "bad request"})
		return
	}
	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}
	err = h.uc.UpdateRewardType(c.Request.Context(), actorID, rewardTypeID, &req)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary Delete reward type by ID
// @Description Delete a reward type by ID
// @Tags rewards
// @Produce json
// @Param id path string true "Reward Type ID"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/rewards/type/{id} [delete]
// @Security ApiKeyAuth
func (h *RewardTypeHandler) DeleteRewardType(c *gin.Context) {
	rewardTypeID, ok := parseUUID(h.logger, c)
	if !ok {
		return
	}
	actorID, ok := parseActorIDFromContext(h.logger, c)
	if !ok {
		return
	}
	err := h.uc.DeleteRewardType(c.Request.Context(), actorID, rewardTypeID)
	if err != nil {
		HandleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}
