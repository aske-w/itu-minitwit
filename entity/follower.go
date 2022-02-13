package entity

import "database/sql"

type Follower struct {
	Who_Id  int64 `db:"who_id" json:"who_id"`
	Whom_Id int64 `db:"whom_id" json:"whom_id"`
}

// TableName returns the database table name of a Product.
func (p Follower) TableName() string {
	return "follower"
}

// PrimaryKey returns the primary key of a Product.
// func (p *Follower) PrimaryKey() string {
// 	return "Follower_id"
// }

// SortBy returns the column name that
// should be used as a fallback for sorting a set of Product.
func (p *Follower) SortBy() string {
	return "who_id"
}

// Scan binds mysql rows to this Product.
func (p *Follower) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Who_Id, &p.Whom_Id)
}

// Products is a list of products. Implements the `Scannable` interface.
type Followers []*Follower

// Scan binds mysql rows to this Categories.
func (ps *Followers) Scan(rows *sql.Rows) (err error) {
	cp := *ps
	for rows.Next() {
		p := new(Follower)
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
