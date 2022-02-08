package entity

// Category represents the categories entity.
// Each product belongs to a category, see `Product.CategoryID` field.
// It implements the `sql.Record` and `sql.Sorted` interfaces.
type Message struct {
	Message_ID int64  `db:"message_id" json:"message_id"`
	Author_ID  int64  `db:"author_id" json:"author_id"`
	Text       string `db:"text" json:"tex"`
	Pub_Data   int64  `db:"pub_date" json:"pub_date"`
	Flagged    int64  `db:"flagged" json:"flagged"`
}
