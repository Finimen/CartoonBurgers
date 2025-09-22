package services

import (
	"CartoonBurgers/models"
	"CartoonBurgers/repositories"
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	GenerateFromPassword([]byte, int) ([]byte, error)
	CompareHashAndPassword([]byte, []byte) error
}

type IRepository interface {
	GetUserByUsername(context.Context, string) (string, error)
	CreateUser(context.Context, models.User, string) error
}

type RegisterHandler struct {
	Hasher     IPasswordHasher
	Repository IRepository
	Logger     *slog.Logger
}

type LoginHandler struct {
	Hasher     IPasswordHasher
	Repository IRepository
	JwtKey     []byte
	Logger     *slog.Logger
}

type BcryptHasher struct {
}

func NewRegisterHandler(hasher IPasswordHasher, repo *repositories.UserRepository, loger *slog.Logger) RegisterHandler {
	var register = RegisterHandler{Hasher: hasher, Repository: repo, Logger: loger}
	return register
}

func NewLoginHandler(hasher IPasswordHasher, repo *repositories.UserRepository, jwtKey []byte, loger *slog.Logger) LoginHandler {
	var login = LoginHandler{Hasher: hasher, Repository: repo, JwtKey: jwtKey, Logger: loger}
	return login
}

func (b *BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (b *BcryptHasher) CompareHashAndPassword(storedPaswsord []byte, userPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(storedPaswsord, userPassword)

	if err != nil {
		return err
	}

	return nil
}
