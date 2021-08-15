package usecase

import (
	"cat-test/internal/domain"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type authUsecase struct {
	repository     domain.AuthRepository
	userRepository domain.UserRepository
}

func (a authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if !a.checkPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}

	token := a.generateNewToken()
	return token, a.repository.CreateToken(ctx, token, user.ID)
}

func (a authUsecase) Logout(ctx context.Context, userID int64) error {
	return a.repository.DestroyToken(ctx, userID)
}

func (a authUsecase) HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 2)
	return string(b), err
}

func (a authUsecase) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a authUsecase) generateNewToken() string {
	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (a authUsecase) Authenticate(ctx context.Context, token string) (domain.User, error) {
	return a.repository.GetUserByToken(ctx, token)
}

func NewAuthUsecase(repository domain.AuthRepository, userRepository domain.UserRepository) domain.AuthUsecase {
	return authUsecase{repository: repository, userRepository: userRepository}
}
