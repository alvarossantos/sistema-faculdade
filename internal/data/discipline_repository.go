package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"

	"github.com/lib/pq"
)

type DisciplineRepository struct {
	DB *sql.DB
}

func (r *DisciplineRepository) GetAll() ([]models.Discipline, error) {
	query := `
		SELECT d.id, d.name, d.code, d.credits, d.workload_hours, d.description,
		       d.department_id, dep.name AS department_name,
		       d.created_at, d.updated_at
		FROM disciplines d
		LEFT JOIN departments dep ON d.department_id = dep.id
		ORDER BY d.name ASC;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar disciplinas: %w", err)
	}
	defer rows.Close()

	var list []models.Discipline

	for rows.Next() {
		var d models.Discipline
		var deptName sql.NullString

		err := rows.Scan(
			&d.ID, &d.Name, &d.Code, &d.Credits, &d.WorkloadHours, &d.Description,
			&d.DepartmentID, &deptName, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear disciplina: %w", err)
		}

		if deptName.Valid {
			d.DepartmentName = deptName.String
		} else {
			d.DepartmentName = "Sem Departamento"
		}
		list = append(list, d)
	}
	return list, nil
}

func (r *DisciplineRepository) Create(d *models.Discipline) (int, error) {
	query := `
		INSERT INTO disciplines (name, code, credits, workload_hours, description, department_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int
	err := r.DB.QueryRow(
		query,
		d.Name, d.Code, d.Credits, d.WorkloadHours, d.Description, d.DepartmentID,
	).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "unique_violation" {
			return 0, fmt.Errorf("código da disciplina já existe: %w", err)
		}
		return 0, fmt.Errorf("erro ao criar disciplina: %w", err)
	}
	return id, nil
}

func (r *DisciplineRepository) GetByID(id int) (*models.Discipline, error) {
	query := `
		SELECT d.id, d.name, d.code, d.credits, d.workload_hours, d.description,
		       d.department_id, dep.name AS department_name,
		       d.created_at, d.updated_at
		FROM disciplines d
		LEFT JOIN departments dep ON d.department_id = dep.id
		WHERE d.id = $1;
	`

	var d models.Discipline
	var deptName sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&d.ID, &d.Name, &d.Code, &d.Credits, &d.WorkloadHours, &d.Description,
		&d.DepartmentID, &deptName, &d.CreatedAt, &d.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar disciplina por ID: %w", err)
	}

	if deptName.Valid {
		d.DepartmentName = deptName.String
	} else {
		d.DepartmentName = "Sem Departamento"
	}

	return &d, nil
}

func (r *DisciplineRepository) Update(d *models.Discipline) error {
	query := `
		UPDATE disciplines
		SET name = $1, code = $2, credits = $3, workload_hours = $4,
		    description = $5, department_id = $6, updated_at = NOW()
		WHERE id = $7
	`

	result, err := r.DB.Exec(
		query,
		d.Name, d.Code, d.Credits, d.WorkloadHours,
		d.Description, d.DepartmentID, d.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar disciplina: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("nenhuma disciplina foi encontrada")
	}
	return nil
}

func (r *DisciplineRepository) Delete(id int) error {
	query := `
		DELETE FROM disciplines
		WHERE id = $1
	`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar disciplina: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("nenhuma disciplina foi encontrada para deletar")
	}
	return nil
}
