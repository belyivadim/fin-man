package models_test

import (
	"testing"

	"github.com/belyivadim/fin-man/models"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
  passwords := []string{"pass", "q!qwG&^233Sdj", "another_one_simple_password"}

  for _, p := range passwords {
    t.Run(p, func(t *testing.T){
      mock_user := models.User{
        Name: "",
        Email: "",
        Password: "",
      }
      err := mock_user.HashPassword(p)
      if err != nil {
        t.Errorf("HashPassword failed with password: %s\n", p)
      }

      err = mock_user.CheckPassword(p)
      if err != nil {
        t.Errorf("CheckPassword failed with password: %s\n", p)
      }
    })
  }
}
