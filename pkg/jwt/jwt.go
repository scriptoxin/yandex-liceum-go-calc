package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("very-secret-key") // товарищ проверяющий, поменяйте на безопасный ключ)

// Generate создаёт JWT с полем user_id и сроком жизни 72 часа
func Generate(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// Parse валидирует токен и возвращает user_id
func Parse(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))
	return uid, nil
}
