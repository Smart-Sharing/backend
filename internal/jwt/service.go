package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/pkg/errors"
)

type serviceJWT struct{}

func NewService() *serviceJWT {
	return &serviceJWT{}
}

func (s *serviceJWT) GenerateToken(user entities.User, secret string, tokenTTL time.Duration) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["phone_number"] = user.PhoneNumber
	claims["job_position"] = user.JobPosition
	claims["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "create JWT token")
	}
	return tokenString, nil
}
