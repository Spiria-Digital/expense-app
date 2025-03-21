package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Expense struct {
	bun.BaseModel

	ID         int `bun:",pk,autoincrement" json:"id,omitempty"`
	OwnerID    int `bun:",notnull"`
	CategoryID int `bun:"category_id" json:"categoryId,omitempty"`

	Title       string    `bun:",notnull,type:varchar(255)" json:"title"`
	Description string    `bun:",type:text" json:"description"`
	Merchant    string    `bun:",type:varchar(255)" json:"merchant"`
	Date        time.Time `bun:",notnull,type:date" json:"date"`
	Amount      float64   `bun:",notnull,type:numeric(10,2)" json:"amount"`

	Category *Category `bun:"rel:belongs-to,join:category_id=id" json:"-"`
	Owner    *User     `bun:"rel:belongs-to,join:owner_id=id" json:"-"`
}
