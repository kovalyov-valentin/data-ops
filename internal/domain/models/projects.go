package models

import "time"

type Projects struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
