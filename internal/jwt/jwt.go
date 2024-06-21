package jwt

import (
	"time"
	"url-shorter/internal/models"

	"github.com/golang-jwt/jwt"
)

func ParseJWT(tokenStr, jwtSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, models.ErrInvalidToken
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", models.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		if exp, ok := (*claims)["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return "", models.ErrTokenExpired
			}
		}
		if login, ok := (*claims)["sub"].(string); ok {
			return login, nil
		}
	}

	return "", models.ErrInvalidToken
}
