package models

import "time"

type Goods struct {
	ID          int       `db:"id"`
	ProjectsID  int       `db:"projects_id"`
	Priority    int       `db:"priority"`
	Removed     bool      `db:"removed"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}
