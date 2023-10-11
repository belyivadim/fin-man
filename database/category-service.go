package database

import "github.com/belyivadim/fin-man/models"

type CategoryService interface {
  GetAllCategories() ([]models.Category, error)
  GetCategoryById(id uint) (*models.Category, error)
  CreateCategoryRecord(category *models.Category) error
  UpdateCategoryRecord(category *models.Category) error
  DeleteCategoryRecordById(id uint) error
}
