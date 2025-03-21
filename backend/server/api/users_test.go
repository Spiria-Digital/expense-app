package api

import (
	"context"
	"encoding/json"
	"github.com/Spiria-Digital/expense-manager/server/middleware"
	"github.com/Spiria-Digital/expense-manager/server/models"
	"github.com/Spiria-Digital/expense-manager/server/service"
	"github.com/Spiria-Digital/expense-manager/server/storage"
	"github.com/Spiria-Digital/expense-manager/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestUserEndpoints(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(storage.GetRootDir(), "expenses.db")
	db, err := storage.NewBunDB(filePath)
	require.NoError(t, err)

	hashedPassword, err := utils.HashPassword("EarthIsFlat#123")
	require.NoError(t, err)
	user := &models.User{
		Email:     "flatearth.1@space.com",
		Password:  hashedPassword,
		FirstName: "Flattie",
		LastName:  "Earther",
	}

	require.NoError(t, service.CreateUser(context.Background(), db, user))

	token, err := middleware.GenerateToken(user.ID)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, service.DeleteUser(context.Background(), db, user.ID))
		require.NoError(t, db.Close())
	})

	t.Run("create and list users", func(t *testing.T) {
		t.Parallel()

		// List users
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("db", db)
		ctx.Request = httptest.NewRequest("GET", "/api/users", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+token)
		ListUsers(ctx)

		require.Equal(t, 200, w.Code)
		var users []map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &users))
		assert.GreaterOrEqual(t, len(users), 1)

		// check if the user is in the list
		exists := false
		for _, u := range users {
			lastName := u["lastName"].(string)
			firstName := u["firstName"].(string)
			if lastName == user.LastName && firstName == user.FirstName {
				exists = true
			}
			// check that password is not returned
			_, ok := u["password"]
			assert.False(t, ok, "password field should not be returned")
		}
		require.True(t, exists, "user not found in the list")
	})
}
