package handlers

import (
	"sistema-faculdade/internal/data"
)

type Handler struct {
	Students    data.StudentRepository
	Teachers    data.TeacherRepository
	Courses     data.CourseRepository
	Departments data.DepartmentRepository
}

func NewHandler(
	s data.StudentRepository,
	t data.TeacherRepository,
	c data.CourseRepository,
	d data.DepartmentRepository,
) *Handler {
	return &Handler{
		Students:    s,
		Teachers:    t,
		Courses:     c,
		Departments: d,
	}
}
