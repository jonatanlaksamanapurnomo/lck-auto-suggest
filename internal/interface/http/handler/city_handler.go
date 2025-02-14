package handler

import (
	"encoding/json"
	"lck-auto-suggest/internal/usecase/city"
	"net/http"
	"strconv"
)

type CityHandler struct {
	service *city.Service
}

func NewCityHandler(service *city.Service) *CityHandler {
	return &CityHandler{
		service: service,
	}
}

func (h *CityHandler) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var lat, lon *float64
	if latStr := r.URL.Query().Get("latitude"); latStr != "" {
		if latVal, err := strconv.ParseFloat(latStr, 64); err == nil {
			lat = &latVal
		}
	}
	if lonStr := r.URL.Query().Get("longitude"); lonStr != "" {
		if lonVal, err := strconv.ParseFloat(lonStr, 64); err == nil {
			lon = &lonVal
		}
	}

	result, err := h.service.GetSuggestions(query, lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
