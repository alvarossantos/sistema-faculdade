package data

import (
	"database/sql"
	"log"
)

type DashboardStats struct {
	Students    int `json:"students"`
	Teachers    int `json:"teachers"`
	Courses     int `json:"courses"`
	Disciplines int `json:"disciplines"`
	Departments int `json:"departments"`
	Semesters   int `json:"semesters"`
}

type DashboardRepository struct {
	DB *sql.DB
}

func (r *DashboardRepository) GetStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	err := r.DB.QueryRow("SELECT COUNT(*) FROM students WHERE active = true").Scan(&stats.Students)
	if err != nil {
		log.Println("Erro ao ao obter contagem de estudantes:", err)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM teachers WHERE active = true").Scan(&stats.Teachers)
	if err != nil {
		log.Println("Erro ao ao obter contagem de professores:", err)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM courses").Scan(&stats.Courses)
	if err != nil {
		log.Println("Erro ao ao obter contagem de cursos:", err)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM disciplines").Scan(&stats.Disciplines)
	if err != nil {
		log.Println("Erro ao ao obter contagem de disciplinas:", err)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM departments").Scan(&stats.Departments)
	if err != nil {
		log.Println("Erro ao ao obter contagem de departamentos:", err)
	}

	err = r.DB.QueryRow("SELECT COUNT(*) FROM academic_semesters").Scan(&stats.Semesters)
	if err != nil {
		log.Println("Erro ao ao obter contagem de semestres:", err)
	}

	return stats, nil
}
