package service

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
)

// ListExpenses returns all expenses for a given user.
// TODO - change the date from the database to return a date in the format "YYYY-MM-DD"
// TODO - refactor ListExpenses and ListExpensesByCategory to use a single function
func ListExpenses(ctx context.Context, db *bun.DB, owner int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := db.NewSelect().Model(&expenses).Where("owner_id = ?", owner).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func ListExpensesByCategory(ctx context.Context, db *bun.DB, owner int, category string) ([]models.Expense, error) {
	var expenses []models.Expense
	err := db.NewSelect().Model(&expenses).Where("owner_id = ? and category_id = ?", owner, category).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func CreateExpense(ctx context.Context, db *bun.DB, expense *models.Expense) error {
	_, err := db.NewInsert().Model(expense).Exec(ctx)
	return err
}

func GetExpense(ctx context.Context, db *bun.DB, id int, owner int) (*models.Expense, error) {
	expense := new(models.Expense)
	err := db.NewSelect().Model(expense).Where("id = ? and owner_id = ?", id, owner).Scan(ctx)
	return expense, err
}

func UpdateExpense(ctx context.Context, db *bun.DB, expense *models.Expense) error {
	_, err := db.NewUpdate().Model(expense).WherePK().Exec(ctx)
	return err
}

func DeleteExpense(ctx context.Context, db *bun.DB, id int, owner int) error {
	_, err := db.NewDelete().Model(&models.Expense{}).Where("id = ? and owner_id = ?", id, owner).Exec(ctx)
	return err
}
