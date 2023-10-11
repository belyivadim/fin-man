package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/belyivadim/fin-man/controllers"
	"github.com/belyivadim/fin-man/models"
	"github.com/gin-gonic/gin"
)


type AuthMockService struct {
  db *[]models.User
}

func NewAuthMockService(db *[]models.User) *AuthMockService {
  return &AuthMockService{
    db: db,
  }
}

func (service *AuthMockService) CreateUserRecord(user *models.User) error {
  *service.db = append(*service.db, *user)
  return nil
}

func (service *AuthMockService) GetUserByEmail(email string) (*models.User, error) {
  for _, u := range *service.db {
    if u.Email == email {
      return &u, nil
    }
  }

  return nil, errors.New("User not found!")
}

func (service *AuthMockService) GetAllUsers() ([]models.User, error) {
  return nil, errors.New("Not Implemented")
}

var mockDb []models.User

func mockEngine() *gin.Engine {
  r := gin.Default()
  v1 := r.Group("/api/v1")

  auth := v1.Group("/auth", func(c *gin.Context) {
    c.Set("AuthService", NewAuthMockService(&mockDb))
    c.Next()
  })

  {
    auth.POST("/signup", controllers.Signup)
    auth.POST("/signin", controllers.Signin)
  }

  return r
}

func TestSignupAndSignIn(t *testing.T) {
  t.Setenv("SECRET_KEY", "mocksecret")
  r := mockEngine()
  var testUsers []models.User

  for i := 0; i < 5; i++ {
    testUsers = append(testUsers, models.User{
      Name: fmt.Sprintf("user%v", i),
      Password: fmt.Sprintf("password%v", i),
      Email: fmt.Sprintf("user%v@mail.com", i),
    })
  }

  for i, u := range testUsers {
    uJson, err := json.Marshal(u)
    if err != nil {
      t.Error("Error marshaling user")
    }

    req, err := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(uJson))
    if err != nil {
      t.Error("Error building request")
    }

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
      t.Errorf("Test /signup failed StatusCode: actual - %v, expected - %v\n", w.Code, http.StatusCreated)
    }

    if len(mockDb) != i + 1 {
      t.Errorf("Test /signup did not actualy write to mockDb: i - %v, len(mockDb) - %v", i, len(mockDb))
    }
  }

  for _, u := range testUsers {
    loginPayload := controllers.LoginPayload {
      Email: u.Email,
      Password: u.Password,
    }

    loginPayloadJson, err := json.Marshal(loginPayload)
    if err != nil {
      t.Error("Error marshaling login payload")
    }

    req, err := http.NewRequest("POST", "/api/v1/auth/signin", bytes.NewBuffer(loginPayloadJson))
    if err != nil {
      t.Error("Error building request")
    }

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
      t.Errorf("Test /signin failed StatusCode: actual - %v, expected - %v\n", w.Code, http.StatusOK)
      var response controllers.ApiError
      err = json.Unmarshal(w.Body.Bytes(), &response)
      if err == nil {
        t.Errorf("Response: %v", response)
      }
    }

    var response controllers.LoginResponse
    err = json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
      t.Errorf("Test /signin failed unmarshaling response")
    }
  }
}
