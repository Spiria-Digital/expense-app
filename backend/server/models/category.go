package models

import "github.com/uptrace/bun"

type Category struct {
	bun.BaseModel

	ID   int    `bun:",pk,autoincrement" json:"id,omitempty"`
	Name string `bun:",unique,notnull" json:"name"`
}
