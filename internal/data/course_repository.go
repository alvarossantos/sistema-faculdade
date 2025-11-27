package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"
)

type CourseRepository struct {
	DB *sql.DB
}

func (r *CourseRepository) GetAll() ([]models.Course, error) {
	query := `SELECT c.id, c.name, c.total_credits_required, c.duration_semesters, c.created_at
	FROM courses c 
	ORDER BY name ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var c models.Course
		err := rows.Scan(
			&c.ID, &c.Name, &c.TotalCreditsRequired, &c.DurationSemesters, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados dos cursos: %w", err)
	}

	return courses, nil
}

func (r *CourseRepository) Create(c *models.Course) (int, error) {
	query := `
		INSERT INTO courses (name, total_credits_required, duration_semesters)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var id int

	err := r.DB.QueryRow(
		query,
		c.Name, c.TotalCreditsRequired, c.DurationSemesters,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("erro ao criar curso: %w", err)
	}

	return id, nil
}
