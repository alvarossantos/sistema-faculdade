package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sistema-faculdade/internal/models"
	"strconv"
)

func (h *Handler) GetAllOffersHandler(w http.ResponseWriter, r *http.Request) {
	list, err := h.Offers.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao buscar ofertas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) CreateOfferHandler(w http.ResponseWriter, r *http.Request) {
	var o models.DisciplineOffer
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	id, err := h.Offers.Create(&o)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao criar oferta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "message": "Oferta criada com sucesso"})
}

func (h *Handler) DeleteOfferHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.Offers.Delete(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Erro ao deletar oferta", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
