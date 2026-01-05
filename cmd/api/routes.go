package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/departments", app.handlers.CreateDepartmentHandler)
	mux.HandleFunc("GET /api/departments", app.handlers.GetAllDepartmentsHandler)

	mux.HandleFunc("POST /api/courses", app.handlers.CreateCourseHandler)
	mux.HandleFunc("GET /api/courses", app.handlers.GetAllCoursesHandler)

	mux.HandleFunc("POST /api/students", app.handlers.CreateStudentHandler)
	mux.HandleFunc("GET /api/students", app.handlers.GetAllStudentsHandler)
	mux.HandleFunc("GET /api/students/{id}", app.handlers.GetStudentByIDHandler)
	mux.HandleFunc("PUT /api/students/{id}", app.handlers.UpdateStudentHandler)
	mux.HandleFunc("DELETE /api/students/{id}", app.handlers.DeleteStudentHandler)
	mux.HandleFunc("PATCH /api/students/{id}/activate", app.handlers.ActivateStudentHandler)

	mux.HandleFunc("POST /api/teachers", app.handlers.CreateTeacherHandler)
	mux.HandleFunc("GET /api/teachers", app.handlers.GetAllTeachersHandler)
	mux.HandleFunc("GET /api/teachers/{id}", app.handlers.GetTeacherByIDHandler)
	mux.HandleFunc("PUT /api/teachers/{id}", app.handlers.UpdateTeacherHandler)
	mux.HandleFunc("DELETE /api/teachers/{id}", app.handlers.DeleteTeacherHandler)
	mux.HandleFunc("PATCH /api/teachers/{id}/activate", app.handlers.ActivateTeacherHandler)

	mux.HandleFunc("POST /api/disciplines", app.handlers.CreateDisciplinesHandler)
	mux.HandleFunc("GET /api/disciplines", app.handlers.GetAllDisciplinesHandler)
	mux.HandleFunc("GET /api/disciplines/{id}", app.handlers.GetDisciplineByIDHandler)
	mux.HandleFunc("PUT /api/disciplines/{id}", app.handlers.UpdateDisciplineHandler)
	mux.HandleFunc("DELETE /api/disciplines/{id}", app.handlers.DeleteDisciplineHandler)

	mux.HandleFunc("POST /api/semesters", app.handlers.CreateSemesterHandler)
	mux.HandleFunc("GET /api/semesters", app.handlers.GetAllSemestersHandler)
	mux.HandleFunc("DELETE /api/semesters/{id}", app.handlers.DeleteSemesterHandler)

	mux.HandleFunc("GET /api/dashboard/stats", app.handlers.GetDashboardStatsHandler)
	// Servidor de arquivos para o frontend
	// Servir CSS
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("ui/static/css"))))

	// Servir JS
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("ui/static/js"))))

	// Servir HTML
	mux.Handle("/", http.FileServer(http.Dir("ui/html")))

	return mux
}
