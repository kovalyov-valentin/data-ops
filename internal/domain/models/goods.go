package models

import "time"

type Goods struct {
	ID          int       `db:"id" json:"id"`
	ProjectsID  int       `db:"projects_id" json:"projects_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Priority    int       `db:"priority" json:"priority"`
	Removed     bool      `db:"removed" json:"removed"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
