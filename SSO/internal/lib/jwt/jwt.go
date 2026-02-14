package jwt

import (
	"authService/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: write test for this func
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["is_admin"] = user.IsAdmin
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
