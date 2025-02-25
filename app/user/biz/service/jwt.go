package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/crazyfrankie/todolist/app/user/config"
)

func GenerateToken(uid int, userAgent string) (string, error) {
	claims := &jwt.MapClaims{
		"user_id":    uid,
		"expire_at":  time.Now().Add(time.Hour * 24).Unix(),
		"issuer":     "github.com/crazyfrankie",
		"issue_at":   time.Now().Unix(),
		"user_agent": userAgent,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConf().JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
