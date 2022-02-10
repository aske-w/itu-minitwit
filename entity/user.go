package entity

import "database/sql"

type User struct {
	ID       int64  `db:"user_id" json:"id"`
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Pw_Hash  string `db:"pw_hash" json:"pw_hash"`
}

// TableName returns the database table name of a Product.
func (p User) TableName() string {
	return "user"
}

// PrimaryKey returns the primary key of a Product.
func (p *User) PrimaryKey() string {
	return "user_id"
}

// SortBy returns the column name that
// should be used as a fallback for sorting a set of Product.
func (p *User) SortBy() string {
	return "id"
}

// Scan binds mysql rows to this Product.
func (p *User) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.ID, &p.Username, &p.Email, &p.Pw_Hash)
}

// Products is a list of products. Implements the `Scannable` interface.
type Users []*User

// Scan binds mysql rows to this Categories.
func (ps *Users) Scan(rows *sql.Rows) (err error) {
	cp := *ps
	for rows.Next() {
		p := new(User)
		if err = p.Scan(rows); err != nil {
			return
		}
		cp = append(cp, p)
	}

	if len(cp) == 0 {
		return sql.ErrNoRows
	}

	*ps = cp

	return rows.Err()
}
