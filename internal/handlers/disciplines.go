package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
	"strconv"
)

func (h *Handler) CreateDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	var d models.Discipline

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Disciplines.Create(&d)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao criar disciplina", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "message": "Criado com sucesso!"})
}

func (h *Handler) GetAllDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.Disciplines.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao buscar disciplinas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) UpdateDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input models.Discipline
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	input.ID = id

	err = h.Disciplines.Update(&input)
	if err != nil {
		if err.Error() == fmt.Sprintf("nenhuma disciplina encontrada com o ID %d", id) {
			http.Error(w, "Disciplina não encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro interno ao atualizar", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((map[string]string{"message": "Disciplina atualizada com sucesso"}))
}

func (h *Handler) DeleteDisciplineHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Disciplines.Delete(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("nenhuma disciplina encontrada com o ID %d", id) {
			http.Error(w, "Disciplina não encontrada", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro interno ao atualizar", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetDisciplineByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	d, err := h.Disciplines.GetByID(id)
	if err != nil {
		http.Error(w, "Erro ao buscar disciplina", http.StatusInternalServerError)
		return
	}
	if d == nil {
		http.Error(w, "Disciplina não encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}
