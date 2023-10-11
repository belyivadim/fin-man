package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belyivadim/fin-man/auth"
	"github.com/belyivadim/fin-man/database"
	"github.com/belyivadim/fin-man/env"
	"github.com/belyivadim/fin-man/models"
)

type LoginPayload struct {
  Email       string `json:"email" bindings:"required"`
  Password    string `json:"password" bindings:"required"`
}

type LoginResponse struct {
  Token         string `json:"token"`
}

func notImplemented(c *gin.Context) {
  c.JSON(http.StatusNotImplemented, gin.H{
    "error": "not implemented",
  })
}


func getAuthService(c *gin.Context) database.AuthService {
  authService, ok := c.Copy().MustGet("AuthService").(database.AuthService)
  if !ok {
    log.Fatal("[auth-controller]: Failed to get authService.")
  }

  return authService
}


// @BasePath /api/v1
//
// Signup godoc
// @Summary signup the user into the system
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object}    models.User
// @Failure 400 {object}    ApiError
// @Failure 500 {object}    ApiError
// @Router /auth/signup [post]
func Signup(c *gin.Context) {
  var user models.User
  err := c.ShouldBindJSON(&user)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusBadRequest, ApiError{Error: "Invalid inputs"})
    c.Abort()
    return
  }

  err = user.HashPassword(user.Password)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error hashing password"})
    c.Abort()
    return
  }

  authService := getAuthService(c)

  err = authService.CreateUserRecord(&user)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error creating user"})
    c.Abort()
    return
  }

  c.JSON(http.StatusCreated, user)
}

// @BasePath /api/v1
//
// Signin godoc
// @Summary signin the user, get jwt token
// @Tags auth
// @Accept json
// @Produce json
// @Param login_payload body LoginPayload true "Login Payload"
// @Success 200 {object}    LoginResponse
// @Failure 401 {object}    ApiError
// @Failure 500 {object}    ApiError
// @Router /auth/signin [post]
func Signin(c *gin.Context) {
  var payload LoginPayload
  err := c.ShouldBindJSON(&payload)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "Error": "Invalid Inputs",
    })
    c.Abort()
    return
  }

  authService := getAuthService(c)
  user, err := authService.GetUserByEmail(payload.Email)
  if err != nil {
    c.JSON(http.StatusUnauthorized, ApiError{Error: "Invalid user credentials"})
    c.Abort()
    return
  }

  err = user.CheckPassword(payload.Password)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusUnauthorized, ApiError{Error: "Invalid user credentials"})
    c.Abort()
    return
  }

  secretKey, err := env.GetSecretKey()
  if err != nil {
    c.JSON(http.StatusInternalServerError, ApiError{
      Error: "Error reading secret key for authentication",
    })
    c.Abort()
    return
  }

  jwtWrapper := auth.JwtWrapper{
    SecretKey: secretKey, 
    Issuer: "AuthService",
    ExpirationMinutes: 30,
  }

  signedToken, err := jwtWrapper.GenerateToken(user.Email)
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error signing token"})
    c.Abort()
    return
  }

  tokenResponse := LoginResponse{
    Token: signedToken,
  }

  c.JSON(http.StatusOK, tokenResponse)
}


// @BasePath /api/v1
//
// GetAllUsers godoc
// @Summary get list of registred users (debug only)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object}    []models.User
// @Failure 500 {object}    ApiError
// @Router /auth/users [get]
func GetAllUsers(c *gin.Context) {
  authService := getAuthService(c)
  users, err := authService.GetAllUsers()
  if err != nil {
    log.Println(err)
    c.JSON(http.StatusInternalServerError, ApiError{Error: "Error fetching users"})
  }

  c.JSON(http.StatusOK, users)
}

