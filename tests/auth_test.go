package tests

import (
	"CartoonBurgers/models"
	"CartoonBurgers/services"
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMocks   func(*MockUserRepository, *MockPasswordHasher)
		expectedCode int
		expectedBody string
	}{
		{
			name: "Succesful regestration",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "password123",
				"email":    "test@gmail.com",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				expectedUser := models.User{
					Username: "validuser",
					Password: "password123",
					Email:    "test@gmail.com",
				}
				mur.On("CreateUser", mock.Anything, expectedUser, "hashed_password").Return(nil)
				mph.On("GenerateFromPassword", []byte("password123"), bcrypt.DefaultCost).Return([]byte("hashed_password"), nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "User registered successfully",
		},
		{
			name: "empty username",
			requestBody: map[string]interface{}{
				"username": "",
				"password": "password123",
				"email":    "test@gmail.com",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "empty password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "",
				"email":    "test@gmail.com",
			},
			setupMocks:   func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "username already exists",
			requestBody: map[string]interface{}{
				"username": "existinguser",
				"password": "password123",
				"email":    "test@gmail.com",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockHasher *MockPasswordHasher) {
				expectedUser := models.User{
					Username: "existinguser",
					Password: "password123",
					Email:    "test@gmail.com",
				}
				mockHasher.On("GenerateFromPassword", mock.Anything, mock.Anything).
					Return([]byte("hashed_password"), nil)
				mockRepo.On("CreateUser", mock.Anything, expectedUser, "hashed_password").
					Return(errors.New("ErrUserExists"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Username already exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(MockUserRepository)
			mockHasher := new(MockPasswordHasher)

			tt.setupMocks(mockRepository, mockHasher)

			handler := services.RegisterHandler{
				Repository: mockRepository,
				Hasher:     mockHasher,
				Logger:     slog.Default(),
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = createTestRequest(http.MethodPost, "/register", tt.requestBody)

			handler.Register(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedBody)

			mockRepository.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}

func TestLoginHandler_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  map[string]interface{}
		setupMocks   func(*MockUserRepository, *MockPasswordHasher)
		expectedCode int
		expectedBody string
		checkToken   bool
	}{
		{
			name: "Successful login",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "correctpassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mur.On("GetUserByUsername", mock.Anything, "validuser").Return(string(hashedPassword), nil)
				mph.On("CompareHashAndPassword", mock.Anything, mock.Anything).Return(nil)
			},
			expectedCode: http.StatusOK,
			checkToken:   true,
		},
		{
			name: "Wrong password",
			requestBody: map[string]interface{}{
				"username": "validuser",
				"password": "wrongpassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
				mur.On("GetUserByUsername", mock.Anything, "validuser").Return(string(hashedPassword), nil)
				mph.On("CompareHashAndPassword", mock.Anything, mock.Anything).
					Return(bcrypt.ErrMismatchedHashAndPassword)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid credits",
		},
		{
			name: "User not found",
			requestBody: map[string]interface{}{
				"username": "nonexistent",
				"password": "anypassword",
			},
			setupMocks: func(mur *MockUserRepository, mph *MockPasswordHasher) {
				mur.On("GetUserByUsername", mock.Anything, "nonexistent").Return("", sql.ErrNoRows)
			},
			expectedCode: http.StatusUnauthorized,
			expectedBody: "Invalid credits",
		},
		{
			name: "Empty username",
			requestBody: map[string]interface{}{
				"username": "",
				"password": "password123",
			},
			setupMocks:   func(mur *MockUserRepository, mph *MockPasswordHasher) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
		{
			name: "Empty password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "",
			},
			setupMocks:   func(mur *MockUserRepository, mph *MockPasswordHasher) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository := new(MockUserRepository)
			mockHasher := new(MockPasswordHasher)

			tt.setupMocks(mockRepository, mockHasher)

			handler := services.LoginHandler{
				Repository: mockRepository,
				Hasher:     mockHasher,
				JwtKey:     []byte("test-secret-key"),
				Logger:     slog.Default(),
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = createTestRequest(http.MethodPost, "/login", tt.requestBody)

			handler.Login(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}

			if tt.checkToken {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				tokenString, exists := response["token"]
				assert.True(t, exists)
				assert.NotEmpty(t, tokenString)

				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte("test-secret-key"), nil
				})

				assert.NoError(t, err)
				assert.True(t, token.Valid)

				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					assert.Equal(t, "validuser", claims["username"])
					assert.NotEmpty(t, claims["exp"])
				}
			}

			mockRepository.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	return redis.NewStringResult(args.String(0), args.Error(1))
}

func (m *MockRedisClient) Exists(key string) *redis.IntCmd {
	args := m.Called(key)
	// ИСПРАВЛЕНИЕ: возвращаем int64 вместо int
	return redis.NewIntResult(args.Get(0).(int64), args.Error(1))
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, expiration)
	return redis.NewStatusResult(args.String(0), args.Error(1))
}

func (m *MockRedisClient) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(key, value, expiration)
	return redis.NewBoolResult(args.Bool(0), args.Error(1))
}

func TestAuthMiddleware_TableDriven(t *testing.T) {
	jwtKey := "test-secret-key"
	logger := slog.Default()

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupRedis     func(*MockRedisClient)
		expectedCode   int
		expectedBody   string
		shouldCallNext bool
	}{
		{
			name: "Valid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/protected", nil)
				token := generateValidToken(jwtKey, "testuser")
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			setupRedis: func(mrc *MockRedisClient) {
				// Токен не в blacklist - возвращаем int64(0)
				mrc.On("Exists", mock.AnythingOfType("string")).Return(int64(0), nil)
			},
			expectedCode:   http.StatusOK,
			expectedBody:   "Authorized",
			shouldCallNext: true,
		},
		{
			name: "Missing token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/protected", nil)
				// Нет заголовка Authorization
				return req
			},
			setupRedis:     func(mrc *MockRedisClient) {},
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   "Missing token",
			shouldCallNext: false,
		},
		{
			name: "Token in blacklist",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/protected", nil)
				token := generateValidToken(jwtKey, "testuser")
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			setupRedis: func(mrc *MockRedisClient) {
				// Токен в blacklist - возвращаем int64(1)
				mrc.On("Exists", mock.AnythingOfType("string")).Return(int64(1), nil)
			},
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   "Token revoked",
			shouldCallNext: false,
		},
		{
			name: "Expired token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/protected", nil)
				token := generateExpiredToken(jwtKey, "testuser")
				req.Header.Set("Authorization", "Bearer "+token)
				return req
			},
			setupRedis: func(mrc *MockRedisClient) {
				mrc.On("Exists", mock.AnythingOfType("string")).Return(int64(0), nil)
			},
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   "Invalid token",
			shouldCallNext: false,
		},
		{
			name: "Invalid token format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/protected", nil)
				req.Header.Set("Authorization", "Bearer invalid-token")
				return req
			},
			setupRedis: func(mrc *MockRedisClient) {
				mrc.On("Exists", mock.AnythingOfType("string")).Return(int64(0), nil)
			},
			expectedCode:   http.StatusUnauthorized,
			expectedBody:   "Invalid token",
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := &MockRedisClient{}
			tt.setupRedis(mockRedis)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = tt.setupRequest()

			c.Set("redis", mockRedis)

			nextCalled := false
			nextHandler := func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"message": "Authorized"})
			}

			middleware := services.AuthMiddleware(jwtKey, logger)

			middleware(c)

			if !c.IsAborted() && tt.shouldCallNext {
				nextHandler(c)
			}

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
			assert.Equal(t, tt.shouldCallNext, nextCalled)

			mockRedis.AssertExpectations(t)
		})
	}
}

func generateValidToken(jwtKey, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(jwtKey))
	return tokenString
}

func generateExpiredToken(jwtKey, username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(jwtKey))
	return tokenString
}
