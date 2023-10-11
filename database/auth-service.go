package database

import "github.com/belyivadim/fin-man/models"

type AuthService interface{
  CreateUserRecord(user *models.User) error;
  GetUserByEmail(email string) (*models.User, error);

  // just for debug
  GetAllUsers() ([]models.User, error)
}


