package tests

import (
	"CartoonBurgers/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (mock *MockUserRepository) CreateUser(ctx context.Context, user models.User, hashedPassword string) error {
	args := mock.Called(ctx, user, hashedPassword)
	return args.Error(0)
}

func (mock *MockUserRepository) GetUserByUsername(ctx context.Context, name string) (string, error) {
	args := mock.Called(ctx, name)
	return args.String(0), args.Error(1)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (mock *MockPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := mock.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (mock *MockPasswordHasher) CompareHashAndPassword(stored []byte, other []byte) error {
	args := mock.Called(stored, other)
	return args.Error(0)
}

func createTestRequest(method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func executeHandler(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}
