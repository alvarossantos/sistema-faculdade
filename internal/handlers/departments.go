package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"

	_ "github.com/lib/pq"
)

func (h *Handler) CreateDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Department

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Departments.Create(&input)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao criar departamento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Departamento criado com sucesso",
		"id":      id,
	})
}

func (h *Handler) GetAllDepartmentsHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.Departments.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno ao buscar departamentos", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
