package auth

import (
  "errors"
  "time"

  jwt "github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key, issuer and expiration time
type JwtWrapper struct {
  SecretKey           string // key used for signing the JWT token
  Issuer              string // Issuer of the JWT token
  ExpirationMinutes   int64  // Number in minutes the JWT token will be valid for
}

// JwtClaim holds Email of the user along with jwt.StandardClaims
type JwtClaim struct {
  Email       string
  jwt.StandardClaims
}

// Generates a JWT token based on the provided email
func (j *JwtWrapper) GenerateToken(email string) (signedToken string, err error) {
  claims := &JwtClaim{
    Email: email,
    StandardClaims: jwt.StandardClaims{
      ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(j.ExpirationMinutes)).Unix(),
      Issuer: j.Issuer,
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedToken, err = token.SignedString([]byte(j.SecretKey))

  return
}

// Validates the JWT token and returns its claims
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
  token, err := jwt.ParseWithClaims(
    signedToken,
    &JwtClaim{},
    func (token *jwt.Token) (interface{}, error) {
      return []byte(j.SecretKey), nil
    },
  )

  if err != nil {
    return
  }

  claims, ok := token.Claims.(*JwtClaim)
  if !ok {
    err = errors.New("Couldn't parse claims")
    return
  }

  if claims.ExpiresAt < time.Now().Local().Unix() {
    err = errors.New("JWT is expired")
    return
  }

  return
}
