package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"

	"github.com/lib/pq"
)

// Estrutura que vai fazer a conexão com o banco
type StudentRepository struct {
	DB *sql.DB
}

// GetAll recupera todos os estudantes ativos do banco de dados.
func (r *StudentRepository) GetAll() ([]models.Student, error) {
	// A query SQL para selecionar os estudantes.
	// Lembre-se de usar parâmetros ($1, $2...) para evitar SQL Injection.
	query := `
		SELECT s.id, s.name, s.email, s.gender, s.date_birth, s.cpf, s.registration_number, s.active,
		s.course_id, c.name as course_name, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN courses c ON s.course_id = c.id
		ORDER BY s.id DESC
	`

	// A função Query executa a consulta no banco de dados.
	// Se você adicionar parâmetros à sua query, passe-os como argumentos adicionais para a função Query.
	// Ex: rows, err := r.DB.Query(query, courseID)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estudantes: %w", err)
	}
	// O defer garante que rows.Close() seja chamado antes da função retornar, liberando a conexão com o banco.
	defer rows.Close()

	// Inicializa uma slice vazia de estudantes que irá armazenar os resultados.
	var students []models.Student

	// Itera sobre cada linha retornada pela consulta.
	for rows.Next() {
		var s models.Student
		var courseName sql.NullString
		// rows.Scan, pega os dados do banco e coloca nas variaveis de s
		err := rows.Scan(
			&s.ID, &s.Name, &s.Email, &s.Gender, &s.DateBirth,
			&s.CPF, &s.RegistrationNumber, &s.Active, &s.CourseID,
			&courseName,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear estudante: %w", err)
		}

		if courseName.Valid {
			s.CourseName = courseName.String
		} else {
			s.CourseName = "Curso não encontrado"
		}

		// Adiciona o estudante recém-escaneado à slice de estudantes.
		students = append(students, s)
	}
	// Após o loop, verifica se ocorreu algum erro durante a iteração.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados dos estudantes: %w", err)
	}

	return students, nil
}

// Create insere um novo estudante no banco de dados.
func (r *StudentRepository) Create(s *models.Student) (int, error) {
	// A query SQL para inserir um novo estudante.
	query := `
		INSERT INTO students (name, email, gender, date_birth, cpf, registration_number, course_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	// Declara a variável que irá receber o ID retornado pelo banco de dados.
	var id int
	// QueryRow é usado aqui porque esperamos que a query retorne exatamente uma linha (o ID do novo estudante).
	// Os argumentos seguintes a 'query' são os valores para os placeholders ($1, $2, ...).
	err := r.DB.QueryRow(
		query,
		s.Name, s.Email, s.Gender, s.DateBirth,
		s.CPF, s.RegistrationNumber, s.CourseID,
	).Scan(&id) // O .Scan atribui o valor da coluna retornada (neste caso, 'id') para a variável 'id'.

	if err != nil {
		// Verifica se o erro é o do tipo Postgres
		if pgErr, ok := err.(*pq.Error); ok {
			// Código 23505 = Unique Violation
			if pgErr.Code == "23505" {
				// Verifica qual constraint violou
				if pgErr.Constraint == "students_cpf_key" {
					return 0, fmt.Errorf("CPF já cadastrado")
				}
				if pgErr.Constraint == "students_email_key" {
					return 0, fmt.Errorf("este email já cadastrado")
				}
				if pgErr.Constraint == "students_registration_number_key" {
					return 0, fmt.Errorf("está matrícula já está cadastrado")
				}
			}
		}
		return 0, fmt.Errorf("erro ao criar estudante: %w", err)
	}

	return id, nil
}

func (r *StudentRepository) Update(s *models.Student) error {

	query := `
		Update students
		SET name = $1, email = $2, gender = $3, date_birth = $4, cpf = $5, registration_number = $6, course_id = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`
	// Executa a query de atualização e escaneia o ID retornado para a variável idReturned.
	result, err := r.DB.Exec(
		query,
		s.Name,
		s.Email,
		s.Gender,
		s.DateBirth,
		s.CPF,
		s.RegistrationNumber,
		s.CourseID,
		s.ID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estudante: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum aluno encontrado com o ID %d", s.ID)
	}

	return nil
}

func (r *StudentRepository) Delete(id int) error {
	// A query SQL para inativar um estudante.
	// Não devemos excluir o aluno, somente defini-lo como desativado
	query := `
		UPDATE students
		SET active = FALSE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	// r.DB.Exec, diferente do QueryRom, o Exec é quando damos um comando no banco
	// mas não esperamos recever dados de volta
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao executar a inativação do estudante com ID %d: %w", id, err)
	}

	// O banco avisa quantas linhas teve alteração, se voltar no 0 (primeira linha)
	// rodou sem erros
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar as linhas afetadas pela inativação: %w", err)
	}

	// Se nenhuma linha foi afetada, significa que o ID do estudante não foi encontrado.
	if rowsAffected == 0 {
		return fmt.Errorf("nenhum aluno encontrado com o ID %d", id)
	}

	return nil
}

func (r *StudentRepository) Activate(id int) error {
	query := `
		UPDATE students
		SET active = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao reativar estudante: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum aluno encontrado com o ID %d", id)
	}

	return nil
}

func (r *StudentRepository) GetByID(id int) (*models.Student, error) {
	query := `
		SELECT s.id, s.name, s.email, s.gender, s.date_birth, s.cpf, s.registration_number, s.active,
		s.course_id, c.name as course_name, s.created_at, s.updated_at
		FROM students s
		LEFT JOIN courses c ON s.course_id = c.id
		WHERE s.id = $1
	`

	var s models.Student
	var courseName sql.NullString

	err := r.DB.QueryRow(query, id).Scan(
		&s.ID, &s.Name, &s.Email, &s.Gender, &s.DateBirth,
		&s.CPF, &s.RegistrationNumber, &s.Active, &s.CourseID,
		&courseName,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar estudante: %w", err)
	}

	if courseName.Valid {
		s.CourseName = courseName.String
	} else {
		s.CourseName = "Curso não encontrado"
	}

	return &s, nil
}
