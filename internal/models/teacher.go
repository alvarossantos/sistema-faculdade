package models

import (
	"time"
)

type Teacher struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	CPF            string    `json:"cpf"`
	Telephone      string    `json:"telephone"`
	Active         bool      `json:"active"`
	DepartmentID   int       `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	DateContract   time.Time `json:"date_contract"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
