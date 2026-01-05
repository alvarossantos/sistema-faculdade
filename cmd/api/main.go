package main

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"os"
	"fmt"
	"strconv"
	"path/filepath"
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

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Erro ao carregar o diretório atual: ", err)
	}

	envPath := filepath.Join(cwd, ".env")
	
	file, err := os.Open(envPath)
	if err != nil {
		log.Println("Erro ao abrir o arquivo .env: ", envPath)
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
	disciplineRepo := data.DisciplineRepository{DB: db}

	myHandlers := handlers.NewHandler(studentRepo, teacherRepo, courseRepo, deptRepo, disciplineRepo)

	app := &application{
		handlers: myHandlers,
	}

	portStr := os.Getenv("PORT")
	var port int
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		} else {
			log.Println("Variável de ambiente PORT inválida, usando porta padrão 8080")
			port = 8080
		}
	}

	log.Printf("Servidor rodando em: http://localhost:%d", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
