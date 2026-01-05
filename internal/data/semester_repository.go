package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"

	"github.com/lib/pq"
)

type SemesterRepository struct {
	DB *sql.DB
}

func (r *SemesterRepository) GetAll() ([]models.AcademicSemester, error) {
	query := `
		SELECT id, year, period
		FROM academic_semesters
		ORDER BY year DESC, period DESC;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar semestres acadêmicos: %w", err)
	}
	defer rows.Close()

	var list []models.AcademicSemester

	for rows.Next() {
		var s models.AcademicSemester
		err := rows.Scan(
			&s.ID, &s.Year, &s.Period,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear semestre acadêmico: %w", err)
		}
		list = append(list, s)
	}
	return list, nil
}

func (r *SemesterRepository) Create(s *models.AcademicSemester) (int, error) {
	query := `
		INSERT INTO academic_semesters (year, period)
		VALUES ($1, $2)
		RETURNING id;
	`

	var id int
	err := r.DB.QueryRow(query, s.Year, s.Period).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("semestre %d.%d já existe", s.Year, s.Period)
		}
		return 0, fmt.Errorf("erro ao criar semestre acadêmico: %w", err)
	}
	return id, nil
}

func (r *SemesterRepository) Delete(id int) error {
	query := `
		DELETE FROM academic_semesters
		WHERE id = $1;
	`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23503" {
			return fmt.Errorf("não é possível excluir: existem ofertas vinculadas a este semestre acadêmico")
		}
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("semestre acadêmico com ID %d não encontrado", id)
	}
	return nil
}

func (r *SemesterRepository) GetByID(id int) (*models.AcademicSemester, error) {
	query := `
		SELECT id, year, period
		FROM academic_semesters
		WHERE id = $1;
	`

	var s models.AcademicSemester
	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.Year, &s.Period)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("semestre acadêmico com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar semestre acadêmico: %w", err)
	}
	return &s, nil
}
