package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"

	"github.com/lib/pq"
)

type TeacherRepository struct {
	DB *sql.DB
}

func (r *TeacherRepository) GetAll() ([]models.Teacher, error) {
	query := `
		SELECT t.id, t.name, t.email, t.cpf, t.telephone, t.active, 
		t.department_id, d.name as department_name, t.date_contract, 
		t.created_at, t.updated_at
		FROM teachers t
		LEFT JOIN departments d ON t.department_id = d.id
		ORDER BY t.id DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar professores: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher

	for rows.Next() {
		var t models.Teacher
		var departmentName sql.NullString
		err := rows.Scan(
			&t.ID, &t.Name, &t.Email, &t.CPF,
			&t.Telephone, &t.Active, &t.DepartmentID,
			&departmentName,
			&t.DateContract, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear professor: %w", err)
		}

		if departmentName.Valid {
			t.DepartmentName = departmentName.String
		} else {
			t.DepartmentName = "Departamento não encontrado"
		}

		teachers = append(teachers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados dos professores: %w", err)
	}

	return teachers, nil
}

func (r *TeacherRepository) Create(t *models.Teacher) (int, error) {
	query := `
		INSERT INTO teachers (name, email, cpf, telephone, department_id, date_contract)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int

	err := r.DB.QueryRow(
		query,
		t.Name, t.Email, t.CPF,
		t.Telephone, t.DepartmentID, t.DateContract,
	).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if pgErr.Constraint == "teachers_cpf_key" {
					return 0, fmt.Errorf("CPF já cadastrado")
				}
				if pgErr.Constraint == "teachers_email_key" {
					return 0, fmt.Errorf("este email já cadastrado")
				}
			}
		}
		return 0, fmt.Errorf("erro ao criar professor: %w", err)
	}

	return id, nil
}

func (r *TeacherRepository) Update(t *models.Teacher) error {
	query := `
		Update teachers
		SET name = $1, email = $2, cpf = $3, telephone = $4, department_id = $5, date_contract = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	result, err := r.DB.Exec(
		query,
		t.Name,
		t.Email,
		t.CPF,
		t.Telephone,
		t.DepartmentID,
		t.DateContract,
		t.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar professor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum professor encontrado com o ID %d", t.ID)
	}

	return nil
}

func (r *TeacherRepository) GetByID(id int) (*models.Teacher, error) {
	query := `
		SELECT t.id, t.name, t.email, t.cpf, t.telephone, t.active, 
		t.department_id, d.name as department_name,t.date_contract, 
		t.created_at, t.updated_at
		FROM teachers t
		LEFT JOIN departments d ON t.department_id = d.id
		WHERE t.id = $1
	`

	var t models.Teacher
	var departmentName sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&t.ID, &t.Name, &t.Email, &t.CPF,
		&t.Telephone, &t.Active, &t.DepartmentID,
		&departmentName,
		&t.DateContract, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar professor: %w", err)
	}

	if departmentName.Valid {
		t.DepartmentName = departmentName.String
	} else {
		t.DepartmentName = "Curso não encontrado"
	}

	return &t, nil
}

func (r *TeacherRepository) Delete(id int) error {
	query := `
		UPDATE teachers
		SET active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao executar a inativação do professor com ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar as linhas afetadas pela inativação: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum professor encontrado com o ID %d", id)
	}

	return nil
}

func (r *TeacherRepository) Activate(id int) error {
	query := `
		UPDATE teachers
		SET active = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao reativar professor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum professor encontrado com o ID %d", id)
	}

	return nil
}
