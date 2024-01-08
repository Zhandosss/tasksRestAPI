package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"restAPI/internal/model"
	"restAPI/internal/repositories"
	"time"
)

const (
	salt       = "dafh434231kdfneajkfn453434nfdknmf"
	signingKey = "aenf9difdnejkjnfdkdlsndknewtnrnewk"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

type AuthService struct {
	rep repositories.Authorization
}

func NewAuthService(rep repositories.Authorization) *AuthService {
	return &AuthService{
		rep: rep,
	}
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) CreateUser(user model.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.rep.CreateUser(user)
}

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	password = generatePasswordHash(password)
	user, err := s.rep.GetUser(login, password)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if password != user.Password {
		return "", fmt.Errorf("%w", ErrWrongPassword)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.ID,
	})
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(inputToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(inputToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("%w", ErrTokenClaims)
	}
	return claims.UserID, nil
}
