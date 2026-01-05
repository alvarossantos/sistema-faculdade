package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
	"strconv"
)

func (h *Handler) CreateSemesterHandler(w http.ResponseWriter, r *http.Request) {
	var input models.AcademicSemester

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if input.Period != 1 && input.Period != 2 {
		http.Error(w, "Período inválido. Deve ser 1 ou 2.", http.StatusBadRequest)
		return
	}

	id, err := h.Semesters.Create(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno ao criar semestre", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "message": "Semestre criado!"})
}

func (h *Handler) GetAllSemestersHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.Semesters.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao buscar semestres", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) DeleteSemesterHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Semesters.Delete(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("semestre não encontrado com o ID %d", id) {
			http.Error(w, "Semestre não encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro interno ao deletar semestre", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
