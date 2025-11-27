package main

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sistema-faculdade/internal/data"
	"sistema-faculdade/internal/handlers"
	"strings"

	_ "github.com/lib/pq"
)

// Conteiner para nossas dependencias
type application struct {
	handlers *handlers.Handler
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}

	}
}

func main() {
	loadEnv()

	connStr := os.Getenv("DB_DSN")
	if connStr == "" {
		log.Fatal("A variável de ambiente DB_DSN não está definida")
	}

	log.Println("Tentando conectar com:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Conectar ao banco de dados
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}

	studentRepo := data.StudentRepository{DB: db}
	teacherRepo := data.TeacherRepository{DB: db}
	courseRepo := data.CourseRepository{DB: db}
	deptRepo := data.DepartmentRepository{DB: db}

	myHandlers := handlers.NewHandler(studentRepo, teacherRepo, courseRepo, deptRepo)

	app := &application{
		handlers: myHandlers,
	}

	log.Println("Servidor pronto! Conectado ao banco.")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
