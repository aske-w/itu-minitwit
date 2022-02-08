package entity

// Category represents the categories entity.
// Each product belongs to a category, see `Product.CategoryID` field.
// It implements the `sql.Record` and `sql.Sorted` interfaces.
type Follower struct {
	Who_Id  int64 `db:"who_id" json:"who_id"`
	Whom_Id int64 `db:"whom_id" json:"whom_id"`
}
