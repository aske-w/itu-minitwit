package entity

import (
	"time"
)

// Category represents the categories entity.
// Each product belongs to a category, see `Product.CategoryID` field.
// It implements the `sql.Record` and `sql.Sorted` interfaces.
type User struct {
	ID       int64  `db:"user_id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Pw_Hash  string `db:"pw_hash" json:"pw_hash"`

	// We could use: sql.NullTime or unix time seconds (as int64),
	// note that the dsn parameter "parseTime=true" is required now in order to fill this field correctly.
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}
