package models

import (
	"time"
)

type Course struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	TotalCreditsRequired int       `json:"total_credits_required"`
	DurationSemesters    int       `json:"duration_semesters"`
	CreatedAt            time.Time `json:"created_at"`
}
