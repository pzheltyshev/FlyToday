package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
)

const (
	WrongParamSearchFlight = "wrong_param_search_flight"
)

func writeJSONError(w http.ResponseWriter, code string, httpCode int, message string) {

	resp, err := json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    code,
		Message: message,
	})

	if err != nil {
		http.Error(w, err.Error(), httpCode)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpCode)
	_, err = w.Write(resp)

	if err != nil {
		log.Printf("Failed to write JSON error: %v", err)
		return
	}
}

func writeJSON(w http.ResponseWriter, data any) {

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("Failed to write JSON error: %v", err)
		return
	}
}

func (f *FlightHandler) Init() {
	http.HandleFunc("GET /api/v1/flight/search", f.searchFlight)
}

type FlightHandler struct {
	service *flight.FlightService
}

func NewFlightHandler(service *flight.FlightService) *FlightHandler {
	return &FlightHandler{
		service: service,
	}
}

func (h *FlightHandler) searchFlight(w http.ResponseWriter, r *http.Request) {

	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse(time.RFC3339, dateStr)

	if err != nil {
		writeJSONError(w, WrongParamSearchFlight, http.StatusBadRequest, "Wrong flight search params")
		return
	}

	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	ctx := r.Context()

	flight := h.service.SearchFlight(ctx, origin, destination, date)

	writeJSON(w, flight)
}
