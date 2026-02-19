package products

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		slog.Error("Error while retrieving products", "Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *handler) FindProductById(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.ParseInt(r.PathValue("id"), 0, 64)
	if err != nil {
		slog.Error("Error while retrieving the product ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product, err := h.service.FindProductById(r.Context(), productId)
	if err != nil {
		slog.Error("Error while retrieving product", "Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
