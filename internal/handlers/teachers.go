package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
	"strconv"

	_ "github.com/lib/pq"
)

func (h *Handler) CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Teacher

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Teachers.Create(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao criar professor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Professor criado com sucesso",
		"id":      id,
	})
}

func (h *Handler) GetAllTeachersHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.Teachers.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno ao buscar professores", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Teachers.Delete(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("nenhum professor encontrado com o ID %d", id) {
			http.Error(w, "Professor não encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro interno ao deletar", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var input models.Teacher
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	input.ID = id

	err = h.Teachers.Update(&input)
	if err != nil {
		if err.Error() == fmt.Sprintf("nenhum professor encontrado com o ID %d", id) {
			http.Error(w, "Professor não encontrado", http.StatusNotFound)
			return
		}
		log.Println("Erro: ", err)
		http.Error(w, "Erro ao atualizar professor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((map[string]string{"message": "Professor atualizado com sucesso"}))
}

func (h *Handler) GetTeacherByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	teacher, err := h.Teachers.GetByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno do servidor", http.StatusBadRequest)
		return
	}

	if teacher == nil {
		http.Error(w, "Professor não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func (h *Handler) ActivateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Teachers.Activate(id)
	if err != nil {
		http.Error(w, "Erro ao ativar professor", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Professor ativado com sucesso"})
}
