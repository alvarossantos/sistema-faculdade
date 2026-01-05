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
	Semesters   data.SemesterRepository
	Dashboard   data.DashboardRepository
}

func NewHandler(
	s data.StudentRepository,
	t data.TeacherRepository,
	c data.CourseRepository,
	d data.DepartmentRepository,
	disc data.DisciplineRepository,
	sem data.SemesterRepository,
	dash data.DashboardRepository,
) *Handler {
	return &Handler{
		Students:    s,
		Teachers:    t,
		Courses:     c,
		Departments: d,
		Disciplines: disc,
		Semesters:   sem,
		Dashboard:   dash,
	}
}
