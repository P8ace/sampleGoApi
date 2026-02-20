package orders

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var data CreateOrder
	err := decoder.Decode(&data)

	if err != nil {
		slog.Error("Error while retrieving request body", "Error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdOrder, err := h.service.PlaceOrder(r.Context(), data)
	if err != nil {
		slog.Error("Error while retrieving products", "Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}
