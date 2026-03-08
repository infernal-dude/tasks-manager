package domain

import "time"

type Task struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Completed   bool      `db:"completed"`
}
