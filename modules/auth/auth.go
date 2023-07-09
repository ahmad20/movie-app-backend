package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthService struct {
	secretKey string
	expiresIn time.Duration
}

type AuthInterface interface {
	GenerateToken(username string) (string, error)
	VerifyToken(tokenString string) (*Claims, error)
}

func NewService(secretKey string, expiresIn time.Duration) AuthInterface {
	return &AuthService{
		secretKey: secretKey,
		expiresIn: expiresIn,
	}
}

func (a *AuthService) GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.expiresIn).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secretKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (a *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	tkn, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return nil, jwt.ErrSignatureInvalid
	}
	claims, ok := tkn.Claims.(*Claims)

	if !ok || !tkn.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
