package db

import (
	"log/slog"

	"github.com/GeorgiiMalishev/ideas-platform/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, logger *slog.Logger) (uuid.UUID, error) {
	err := db.AutoMigrate(
		&models.User{},
		&models.BannedUser{},
		&models.Role{},
		&models.CoffeeShop{},
		&models.WorkerCoffeeShop{},
		&models.Category{},
		&models.Idea{},
		&models.IdeaLike{},
		&models.IdeaComment{},
		&models.IdeaStatus{},
		&models.Reward{},
		&models.RewardType{},
		&models.OTP{},
		&models.UserRefreshToken{},
	)
	if err != nil {
		return uuid.Nil, err
	}

	statuses := []string{"Создана", "В работе", "Реализована", "Отклонена"}
	for _, title := range statuses {
		err := db.FirstOrCreate(&models.IdeaStatus{Title: title}, "title = ?", title).Error
		if err != nil {
			logger.Error("Failed to seed status:", slog.String("status", title), slog.String("error", err.Error()))
		}
	}

	adminRole := models.Role{
		Name: "admin",
	}
	if err := db.FirstOrCreate(&adminRole, "name = ?", "admin").Error; err != nil {
		return uuid.Nil, err
	}

	return adminRole.ID, nil
}
