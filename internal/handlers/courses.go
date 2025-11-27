package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
)

func (h *Handler) CreateCourseHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Course

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Courses.Create(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao criar curso", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Curso criado com sucesso",
		"id":      id,
	})
}

func (h *Handler) GetAllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	courses, err := h.Courses.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno ao buscar cursos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
