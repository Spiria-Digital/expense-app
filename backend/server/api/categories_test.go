package api

import (
	"bytes"
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

func TestCreateListCategories(t *testing.T) {
	t.Parallel()

	filePath := filepath.Join(storage.GetRootDir(), "expenses.db")
	db, err := storage.NewBunDB(filePath)
	require.NoError(t, err)

	hashedPassword, err := utils.HashPassword("secret123")
	require.NoError(t, err)
	user := &models.User{
		Email:     "category.creator@test.com",
		Password:  hashedPassword,
		FirstName: "John",
		LastName:  "Doe",
	}

	require.NoError(t, service.CreateUser(context.Background(), db, user))

	token, err := middleware.GenerateToken(user.ID)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, service.DeleteUser(context.Background(), db, user.ID))
		require.NoError(t, db.Close())
	})

	t.Run("test create and list", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		payload, _ := json.Marshal(map[string]string{
			"name": "Test Groceries",
		})
		ctx.Request = httptest.NewRequest("POST", "/api/categories", bytes.NewBuffer(payload))
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		ctx.Set("db", db)

		CreateCategory(ctx)

		assert.Equal(t, 201, w.Code)
		var category models.Category
		err = json.Unmarshal(w.Body.Bytes(), &category)
		assert.Equal(t, category.Name, "Test Groceries")
		assert.Greater(t, category.ID, 0)

		defer func() {
			require.NoError(t, service.DeleteCategory(context.Background(), db, &category))
		}()

		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/categories", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		ctx.Set("db", db)

		ListCategories(ctx)

		assert.Equal(t, 200, w.Code)
		var response []models.Category
		err = json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})
}
