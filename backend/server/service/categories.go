package service

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
)

// GetCategories returns all expenses for a given user.
func GetCategories(ctx context.Context, db *bun.DB) ([]models.Category, error) {
	var categories []models.Category
	err := db.NewSelect().Model(&categories).OrderExpr("name").Scan(ctx)
	return categories, err
}

func CreateCategory(ctx context.Context, db *bun.DB, category *models.Category) error {
	_, err := db.NewInsert().Model(category).Returning("id").Exec(ctx)
	return err
}

func UpdateCategory(ctx context.Context, db *bun.DB, category *models.Category) error {
	_, err := db.NewUpdate().Model(category).Column("name").Where("id = ?", category.ID).Exec(ctx)
	return err
}

func DeleteCategory(ctx context.Context, db *bun.DB, category *models.Category) error {
	_, err := db.NewDelete().Model(category).Where("id = ?", category.ID).Exec(ctx)
	return err
}

func GetCategory(ctx context.Context, db *bun.DB, id int) (*models.Category, error) {
	category := new(models.Category)
	err := db.NewSelect().Model(category).Where("id = ?", id).Scan(ctx)
	return category, err
}
