package handlers

import (
	"sistema-faculdade/internal/data"
)

type Handler struct {
	Students    data.StudentRepository
	Teachers    data.TeacherRepository
	Courses     data.CourseRepository
	Departments data.DepartmentRepository
	Disciplines data.DisciplineRepository
}

func NewHandler(
	s data.StudentRepository,
	t data.TeacherRepository,
	c data.CourseRepository,
	d data.DepartmentRepository,
	disc data.DisciplineRepository,
) *Handler {
	return &Handler{
		Students:    s,
		Teachers:    t,
		Courses:     c,
		Departments: d,
		Disciplines: disc,
	}
}
