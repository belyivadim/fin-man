package database

import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

func NewGormDB(dns string) (*gorm.DB, error) {
  return gorm.Open(postgres.Open(dns), &gorm.Config{})
}
