package database

import (
	"gorm.io/gorm"

	"github.com/belyivadim/fin-man/models"
)

type AuthGormService struct {
  db *gorm.DB
}

func NewAuthGormService(db *gorm.DB) *AuthGormService {
  return &AuthGormService{
    db: db,
  }
}

func (service *AuthGormService) CreateUserRecord(user *models.User) error {
  return service.db.Create(&user).Error
}

func (service *AuthGormService) GetUserByEmail(email string) (*models.User, error) {
  var user models.User
  result := service.db.Where("email = ?", email).First(&user)
  return &user, result.Error
}

func (service *AuthGormService) GetAllUsers() ([]models.User, error) {
  var users []models.User
  result := service.db.Find(&users)
  if result.Error != nil {
    return nil, result.Error
  }

  return users, nil
}
