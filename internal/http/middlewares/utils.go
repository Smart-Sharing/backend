package middlewares

import (
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/golang-jwt/jwt/v5"
)

func extractClaims(tokenStr, secret string) (jwt.MapClaims, bool) {
	hmacSecret := []byte(secret)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	}
	return nil, false
}

func tokenExpire(tokenExp int64) bool {
	return time.Now().Unix() >= tokenExp
}

func hasUserPermission(userJob, requiredJob entities.UserJob) bool {
	switch requiredJob {
	case entities.Admin:
		return userJob == entities.Admin

	case entities.Worker:
		return userJob == entities.Worker || userJob == entities.Admin
	}
	return false
}
