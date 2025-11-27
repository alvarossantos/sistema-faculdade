package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"
)

type DepartmentRepository struct {
	DB *sql.DB
}

func (r *DepartmentRepository) GetAll() ([]models.Department, error) {
	query := `
		SELECT d.id, d.name, d.abbreviation, d.created_at
		FROM departments d
		ORDER BY name ASC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar departamentos: %w", err)
	}
	defer rows.Close()

	var departments []models.Department

	for rows.Next() {
		var d models.Department
		err := rows.Scan(
			&d.ID, &d.Name, &d.Abbreviation, &d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear departamento: %w", err)
		}

		departments = append(departments, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados dos departamentos: %w", err)
	}

	return departments, nil
}

func (r *DepartmentRepository) Create(d *models.Department) (int, error) {
	query := `
		INSERT INTO departments (name, abbreviation)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int

	err := r.DB.QueryRow(
		query,
		d.Name, d.Abbreviation,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar departamento: %w", err)
	}

	return id, nil
}
