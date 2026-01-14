package data

import (
	"database/sql"
	"fmt"
	"sistema-faculdade/internal/models"

	"github.com/lib/pq"
)

type OfferRepository struct {
	DB *sql.DB
}

func (r *OfferRepository) GetAll() ([]models.DisciplineOffer, error) {
	query := `
		SELECT o.id, o.discipline_id, d.name, o.teacher_id, t.name,
		o.semester_id, s.year, s.period, o.schedule, o.class_code
		FROM discipline_offers o
		JOIN disciplines d ON o.discipline_id = d.id
		JOIN teachers t ON o.teacher_id = t.id
		JOIN academic_semesters s ON o.semester_id = s.id
		ORDER BY s.year DESC, s.period DESC, d.name ASC;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ofertas: %w", err)
	}
	defer rows.Close()

	var list []models.DisciplineOffer

	for rows.Next() {
		var o models.DisciplineOffer
		var year, period int
		err := rows.Scan(
			&o.ID, &o.DisciplineID, &o.DisciplineName,
			&o.TeacherID, &o.TeacherName,
			&o.SemeterID, &year, &period,
			&o.Schedule, &o.ClassCode,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear oferta: %w", err)
		}

		o.SetSemesterLabel(year, period)
		list = append(list, o)
	}

	return list, nil
}

func (r *OfferRepository) Create(o *models.DisciplineOffer) (int, error) {
	query := `
		INSERT INTO discipline_offers (discipline_id, semester_id, teacher_id, schedule, class_code)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var id int
	err := r.DB.QueryRow(
		query,
		o.DisciplineID, o.SemeterID, o.TeacherID, o.Schedule, o.ClassCode,
	).Scan(&id)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("oferta já existe para a disciplina %d no semestre %d", o.DisciplineID, o.SemeterID)
		}
		return 0, fmt.Errorf("erro ao criar oferta: %w", err)
	}
	return id, nil
}

func (r *OfferRepository) Delete(id int) error {
	query := `DELETE FROM discipline_offers WHERE id = $1;`
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar oferta: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("nenhuma oferta foi encontrada para deletar")
	}
	return nil
}
