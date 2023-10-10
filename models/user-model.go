package models

import (
  "golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
)

type User struct {
  gorm.Model
  Name      string    `json:"name" binding:"required"`
  Email     string    `json:"email" binding:"required" gorm:"unique"`
  Password  string    `json:"password" binding:"required"`
}


// Takes a password string as a parameter and encrypts it using bcrypt,
// saves the result into user.Password field
func (user *User) HashPassword(password string) error {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  if err != nil {
    return err
  }
  user.Password = string(bytes)
  return nil
}

// Checks the providedPassword for against the user.Password field
func (user *User) CheckPassword(providedPassword string) error {
  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
  return err
}
