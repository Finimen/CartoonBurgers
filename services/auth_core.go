package services

import (
	"CartoonBurgers/repositories"

	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	GenerateFromPassword([]byte, int) ([]byte, error)
	CompareHashAndPassword([]byte, []byte) error
}

type RegisterHandler struct {
	Hasher     IPasswordHasher
	Repository repositories.UserRepository
}

type LoginHandler struct {
	Hasher     IPasswordHasher
	Repository repositories.UserRepository
	JwtKey     []byte
}

type BcryptHasher struct {
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
