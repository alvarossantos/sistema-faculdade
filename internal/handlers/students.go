package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handler) CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	// Variavel para receber o JSON do corpo da requisição
	var input models.Student

	// Decodificador JSON: Le o Body e converte para Struct
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Chama o banco
	id, err := h.Students.Create(&input)
	if err != nil {
		log.Println(err)

		msg := err.Error()
		if strings.Contains(msg, "CPF já cadastrado") ||
			strings.Contains(msg, "email já cadastrado") ||
			strings.Contains(msg, "matrpicula já está cadastrado") {
			http.Error(w, msg, http.StatusConflict)
			return
		}
		http.Error(w, "Erro interno ao criar estudante", http.StatusInternalServerError)
		return
	}

	// Responde para o cliente
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Estudante criado com sucesso",
		"id":      id,
	})
}

func (h *Handler) GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Chama o banco
	list, err := h.Students.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno ao buscar alunos", http.StatusBadRequest)
		return
	}

	// Cabeçalho com JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	// Precisamos pegar o ID da url
	// r.PathValue pega o valor onde definimos {id} na rota
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Chama o banco
	err = h.Students.Delete(id)
	if err != nil {
		// Se nenhum aluno foi encontrado retorna 404 (Not Found)
		if err.Error() == fmt.Sprintf("nenhum aluno encontrado com o ID %d", id) {
			http.Error(w, "Aluno não encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro interno ao deletar", http.StatusBadRequest)
		return
	}

	// Responde ao cliente
	// Status 204 No Content, padrão para Deletes que dão certo
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Ler o JSON pelo usuario
	var input models.Student
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Erro ao ler JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Garantir que o ID da struct seja o da URL
	input.ID = id

	// Chama o banco
	err = h.Students.Update(&input)
	if err != nil {
		if err.Error() == fmt.Sprintf("nenhum aluno encontrado com o ID %d", id) {
			http.Error(w, "Aluno não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "Erro interno ao atualizar", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode((map[string]string{"message": "Aluno atualizado com sucesso"}))
}

func (h *Handler) GetStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	student, err := h.Students.GetByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro interno de servidor", http.StatusBadRequest)
		return
	}

	if student == nil {
		http.Error(w, "Aluno não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *Handler) ActivateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Students.Activate(id)
	if err != nil {
		http.Error(w, "Erro ao reativar aluno", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Aluno reativado com sucesso"})
}
