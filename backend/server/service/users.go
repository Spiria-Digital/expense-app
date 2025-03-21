package service

import (
	"context"
	"time"

	"github.com/uptrace/bun"

	"github.com/Spiria-Digital/expense-manager/server/models"
)

var (
	SQLTimeoutDuration = time.Second * 10
)

func CreateUser(ctx context.Context, db *bun.DB, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	_, err := db.NewInsert().Model(user).Returning("id").Exec(ctx)
	return err
}

func GetUserById(ctx context.Context, db *bun.DB, id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	user := new(models.User)
	err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	return user, err
}

func GetUserByEmail(ctx context.Context, db *bun.DB, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	user := new(models.User)
	err := db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	return user, err
}

func UpdateUser(ctx context.Context, db *bun.DB, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	_, err := db.NewUpdate().Model(&user).
		WherePK().Exec(ctx)
	return err
}

func DeleteUser(ctx context.Context, db *bun.DB, id int) error {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	_, err := db.NewDelete().Model(&models.User{}).Where("id = ?", id).Exec(ctx)
	return err
}

// ListUsers returns a list of first 100 users.
func ListUsers(ctx context.Context, db *bun.DB) ([]models.OutgoingUser, error) {
	ctx, cancel := context.WithTimeout(ctx, SQLTimeoutDuration)
	defer cancel()

	var users []models.OutgoingUser
	err := db.NewSelect().
		ModelTableExpr("users as u1").
		Model(&users).
		ColumnExpr("u1.id, u1.first_name, u1.last_name, u1.is_admin").
		Order("last_name ASC").
		Limit(100).
		Scan(ctx)
	return users, err
}
