package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) GetDashboardStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Dashboard.GetStats()
	if err != nil {
		log.Println("Erro no dashboard:", err)
		http.Error(w, "Erro ao carregar estatísticas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
