package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Spiria-Digital/expense-manager/server/models"
	"github.com/Spiria-Digital/expense-manager/server/service"
	"github.com/Spiria-Digital/expense-manager/server/storage"
	"github.com/Spiria-Digital/expense-manager/server/utils"
)

var (
	invalidRequestMsg     = `invalid request body`
	invalidCredentialsMsg = `invalid email or password`
)

func TestUserLogin(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(storage.GetRootDir(), "expenses.db")
	db, err := storage.NewBunDB(filePath)
	require.NoError(t, err)

	hashedPassword, err := utils.HashPassword("secret123")
	require.NoError(t, err)
	user := &models.User{
		Email:     "john.doe@test.com",
		Password:  hashedPassword,
		FirstName: "John",
		LastName:  "Doe",
	}

	require.NoError(t, service.CreateUser(context.Background(), db, user))

	t.Cleanup(func() {
		require.NoError(t, service.DeleteUser(context.Background(), db, user.ID))
		require.NoError(t, db.Close())
	})

	testCases := map[string]struct {
		payload            map[string]interface{}
		partialResponse    string
		expectedStatusCode int
	}{
		"invalid payload": {
			payload:            map[string]interface{}{},
			partialResponse:    invalidRequestMsg,
			expectedStatusCode: 400,
		},
		"missing email": {
			payload: map[string]interface{}{
				"password": "password",
			},
			partialResponse:    invalidRequestMsg,
			expectedStatusCode: 400,
		},
		"missing password": {
			payload: map[string]interface{}{
				"email": "email@test.com",
			},
			partialResponse:    invalidRequestMsg,
			expectedStatusCode: 400,
		},
		"missing user": {
			payload: map[string]interface{}{
				"email":    "email@test.com",
				"password": "password",
			},
			partialResponse:    invalidCredentialsMsg,
			expectedStatusCode: 401,
		},
		"valid user but invalid login": {
			payload: map[string]interface{}{
				"email":    "john.doe@test.com",
				"password": "password",
			},
			partialResponse:    invalidCredentialsMsg,
			expectedStatusCode: 401,
		},
		"valid user and credentials": {
			payload: map[string]interface{}{
				"email":    "john.doe@test.com",
				"password": "secret123",
			},
			partialResponse:    "token",
			expectedStatusCode: 200,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			payload, _ := json.Marshal(tc.payload)
			ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(payload))

			ctx.Set("db", db)

			UserLogin(ctx)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.partialResponse)
		})
	}
}

func TestUserRegistration(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(storage.GetRootDir(), "expenses.db")
	db, err := storage.NewBunDB(filePath)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	testCases := map[string]struct {
		payload            map[string]interface{}
		partialResponse    string
		expectedStatusCode int
	}{
		"invalid payload": {
			payload:            map[string]interface{}{},
			partialResponse:    invalidRequestMsg,
			expectedStatusCode: 400,
		},
		"missing email": {
			payload: map[string]interface{}{
				"password":  "password",
				"firstName": "John",
				"lastName":  "Doe",
			},
			partialResponse:    `Field validation for 'Email' failed on the 'required' tag`,
			expectedStatusCode: 400,
		},
		"missing password": {
			payload: map[string]interface{}{
				"email":     "user1@company.net",
				"firstName": "John",
				"lastName":  "Doe",
			},
			partialResponse:    `Field validation for 'Password' failed on the 'required' tag`,
			expectedStatusCode: 400,
		},
		"missing required fields - first name": {
			payload: map[string]interface{}{
				"email":    "user1@company.net",
				"password": "secret123",
				"lastName": "Doe",
			},
			partialResponse:    `Field validation for 'FirstName' failed on the 'required' tag`,
			expectedStatusCode: 400,
		},
		"missing required fields - last name": {
			payload: map[string]interface{}{
				"email":     "user1@company.net",
				"password":  "secret123",
				"firstName": "John",
			},
			partialResponse:    `Field validation for 'FirstName' failed on the 'required' tag`,
			expectedStatusCode: 400,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			payload, _ := json.Marshal(tc.payload)
			ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(payload))

			ctx.Set("db", db)

			UserRegistration(ctx)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.partialResponse)
		})
	}

	t.Run("valid user registration", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		email := "user1@company.io"
		payload, _ := json.Marshal(map[string]interface{}{
			"email":     email,
			"password":  "secret123",
			"firstName": "John",
			"lastName":  "Doe",
		})
		ctx.Request = httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(payload))

		ctx.Set("db", db)

		UserRegistration(ctx)

		require.Equal(t, 201, w.Code)

		// select the user and validate that password is hashed
		user, err := service.GetUserByEmail(context.Background(), db, email)
		require.NoError(t, err)
		assert.NotEqual(t, "secret123", user.Password)

		// cleanup
		require.NoError(t, service.DeleteUser(context.Background(), db, user.ID))
	})
}
