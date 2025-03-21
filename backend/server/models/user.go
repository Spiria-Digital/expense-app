package models

import (
	"fmt"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel

	ID        int    `bun:",pk,autoincrement" json:"id"`
	Email     string `bun:",unique,notnull" json:"email" binding:"required"`
	Password  string `bun:",type:varchar(255),notnull" binding:"required"`
	FirstName string `bun:",notnull" json:"first_name" binding:"required"`
	LastName  string `bun:",notnull" json:"last_name" binding:"required"`
	IsAdmin   bool   `bun:",notnull,default:false"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

type OutgoingUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsAdmin   bool   `json:"isAdmin"`
}
