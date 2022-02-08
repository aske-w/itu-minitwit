package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	Conn *sql.DB
}

var ErrNoRows = sql.ErrNoRows

func ConnectSqlite() (*SQLite, error) {
	conn, err := sql.Open("sqlite3", "./db.db")

	if err != nil {
		return nil, err
	}

	return &SQLite{
		Conn: conn,
	}, nil

}

func (db *SQLite) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, err := db.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if scannable, ok := dest.(Scannable); ok {
		return scannable.Scan(rows)
	}

	if !rows.Next() {
		return ErrNoRows
	}
	return rows.Scan(dest)

	/* Uncomment this and pass a slice if u want to see reflection powers <3
	v, ok := dest.(reflect.Value)
	if !ok {
		v = reflect.Indirect(reflect.ValueOf(dest))
	}
	sliceTyp := v.Type()
	if sliceTyp.Kind() != reflect.Slice {
		sliceTyp = reflect.SliceOf(sliceTyp)
	}
	sliceElementTyp := deref(sliceTyp.Elem())
	for rows.Next() {
		obj := reflect.New(sliceElementTyp)
		obj.Interface().(Scannable).Scan(rows)
		if err != nil {
			return err
		}
		v.Set(reflect.Append(v, reflect.Indirect(obj)))
	}
	*/
}

// Scannable for go structs to bind their fields.
type Scannable interface {
	Scan(*sql.Rows) error
}
