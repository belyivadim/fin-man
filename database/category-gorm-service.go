package database

import (
	"github.com/belyivadim/fin-man/models"
	"gorm.io/gorm"
)

type CategoryGormService struct {
  db *gorm.DB
}

type CategoryNotFoundError struct{}

func (e *CategoryNotFoundError) Error() string {
  return "Category is not found"
}

func NewCategoryGormService(db *gorm.DB) *CategoryGormService {
  return &CategoryGormService{
    db: db,
  }
}

func (service *CategoryGormService) GetAllCategories() ([]models.Category, error) {
  var categories []models.Category
  result := service.db.Find(&categories)
  return categories, result.Error
}

func (service *CategoryGormService) GetCategoryById(id uint) (*models.Category, error) {
  var category models.Category
  result := service.db.Where("id = ?", id).First(&category)
  return &category, result.Error
}

func (service *CategoryGormService) CreateCategoryRecord(category *models.Category) error {
  return service.db.Create(category).Error
}

func (service *CategoryGormService) UpdateCategoryRecord(category *models.Category) error {
  found_category, err := service.GetCategoryById(category.ID)
  if err != nil {
    return &CategoryNotFoundError{}
  }

  found_category.Name = category.Name
  return service.db.Model(&found_category).Updates(&found_category).Error
}

func (service *CategoryGormService) DeleteCategoryRecordById(id uint) error {
  category := models.Category{}
  category.ID = id
  return service.db.Delete(&category).Error
}

