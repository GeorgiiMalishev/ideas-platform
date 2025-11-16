package handlers

import (
	"log/slog"
	"net/http"

	"github.com/GeorgiiMalishev/ideas-platform/internal/dto"
	"github.com/GeorgiiMalishev/ideas-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	uc     usecase.AuthUsecase
	logger *slog.Logger
}

func NewAuthHandler(uc usecase.AuthUsecase, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *AuthHandler) GetOTP(c *gin.Context) {
	phone := c.Param("phone")
	err := h.uc.GetOTP(phone)
	if err != nil {
		handleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "bad request"})
		return
	}

	authResp, err := h.uc.VerifyOTP(&req)
	if err != nil {
		handleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, authResp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshReq dto.RefreshRequest
	err := c.ShouldBindJSON(&refreshReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "bad request"})
		return
	}

	authResp, err := h.uc.Refresh(refreshReq.RefreshToken)
	if err != nil {
		handleAppErrors(err, h.logger, c)
		return
	}

	c.JSON(http.StatusOK, authResp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var logoutReq dto.LogoutRequest
	err := c.ShouldBindJSON(&logoutReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Message: "bad request"})
		return
	}

	err = h.uc.Logout(logoutReq.RefreshToken)
	if err != nil {
		handleAppErrors(err, h.logger, c)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) LogoutEverywhere(c *gin.Context) {
	userIDAny, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "user not authorized"})
		return
	}

	userID, ok := userIDAny.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "internal server error"})
		return
	}

	err := h.uc.LogoutEverywhere(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
