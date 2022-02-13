package entity

import (
	"database/sql"
)

// adapted from: https://github.com/kataras/iris/blob/4899fe95f47dfb8ce1f309c5ba267a7ecd0e08ba/_examples/database/mysql/entity/product.go
type Message struct {
	Message_id int    `db:"message_id" json:"message_id"`
	Author_id  int    `db:"author_id" json:"author_id"`
	Text       string `db:"text" json:"text"`
	Pub_date   int    `db:"pub_date" json:"pub_date"`
	Flagged    int    `db:"flagged" json:"flagged"`
}

// TableName returns the database table name of a Product.
func (p Message) TableName() string {
	return "message"
}

// PrimaryKey returns the primary key of a Product.
func (p *Message) PrimaryKey() string {
	return "message_id"
}

// SortBy returns the column name that
// should be used as a fallback for sorting a set of Product.
func (p *Message) SortBy() string {
	return "pub_date"
}

// Scan binds mysql rows to this Product.
func (p *Message) Scan(rows *sql.Rows) error {
	return rows.Scan(&p.Message_id, &p.Author_id, &p.Text, &p.Pub_date, &p.Flagged)
}

// Products is a list of products. Implements the `Scannable` interface.
type Messages []*Message

// Scan binds mysql rows to this Categories.
func (ps *Messages) Scan(rows *sql.Rows) (err error) {
	cp := *ps
	for rows.Next() {
		p := new(Message)
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
