package models

type Product struct {
	ID          int64  `db:"id"`
	Category    string `db:"category"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Count       int64  `db:"count"`
}
