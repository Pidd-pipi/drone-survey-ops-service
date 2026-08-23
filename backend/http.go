package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type statusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewRouter(service *SurveyService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/missions", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, service.Missions()) })
	mux.HandleFunc("POST /api/missions/status", func(w http.ResponseWriter, r *http.Request) {
		var req statusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ID) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id and status JSON are required"})
			return
		}
		mission, err := service.ChangeStatus(req.ID, req.Status)
		if errors.Is(err, ErrMissionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, mission)
	})
	return withStatic(mux)
}
