package models

import (
	"time"
)

// Student é a estrutura de um aluno no formato do banco de dados
// As "tags" json ajuda a quando transformamos em api
type Student struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Email              *string   `json:"email"`
	Gender             string    `json:"gender"`
	DateBirth          time.Time `json:"date_birth"`
	CPF                string    `json:"cpf"`
	RegistrationNumber string    `json:"registration_number"`
	Active             bool      `json:"active"`
	CourseID           int       `json:"course_id"`
	CourseName         string    `json:"course_name"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
