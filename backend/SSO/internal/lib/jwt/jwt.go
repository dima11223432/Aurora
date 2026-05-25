// Package jwt provides JWT token generation for authenticated users.
package jwt

import (
	"authService/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: write test for this func
// NewToken generates a signed JWT token for the given user and app with the specified duration.
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["telegram_id"] = user.Telegram_id
	claims["is_admin"] = user.Is_admin
	claims["username"] = user.Username
	claims["first_name"] = user.First_name
	claims["last_name"] = user.Last_name
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
