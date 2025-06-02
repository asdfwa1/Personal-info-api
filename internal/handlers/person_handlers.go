package handlers

import (
	"Test_Task_EffMob/internal/models"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

type PersonService interface {
	CreatePerson(ctx context.Context, p *models.Person) error
	GetPersons(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error)
	UpdatePerson(ctx context.Context, id uint32, p *models.Person) error
	DeletePerson(ctx context.Context, id uint32) error
}

type PersonHandler struct {
	Service PersonService
}

func NewPersonHandler(service PersonService) *PersonHandler {
	return &PersonHandler{
		Service: service,
	}
}

func (ph *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.DebugContext(r.Context(), "Invalid input", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ph.Service.CreatePerson(r.Context(), &p); err != nil {
		slog.DebugContext(r.Context(), "Failed to create person", "error", err)
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (ph *PersonHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterParams{
		Name:    r.URL.Query().Get("name"),
		Surname: r.URL.Query().Get("surname"),
	}

	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		} else {
			slog.DebugContext(r.Context(), "Invalid limit parameter", "error", err)
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	offset := 0
	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		} else {
			slog.DebugContext(r.Context(), "Invalid offset parameter", "error", err)
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	pagination := models.PaginationParams{
		Limit:  limit,
		OffSet: offset,
	}

	persons, err := ph.Service.GetPersons(r.Context(), filter, pagination)
	if err != nil {
		slog.DebugContext(r.Context(), "Failed getting list person", "error", err)
		http.Error(w, "Failed getting list person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(persons)
}

func (ph *PersonHandler) Update(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		slog.DebugContext(r.Context(), "Failed parse id to number (in Update)", "error", err)
		http.Error(w, "Failed parse id to number", http.StatusBadRequest)
		return
	}
	id := uint32(id64)

	var p models.Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.DebugContext(r.Context(), "Invalid input", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := ph.Service.UpdatePerson(r.Context(), id, &p); err != nil {
		slog.DebugContext(r.Context(), "Failed update person", "error", err)
		http.Error(w, "Failed update person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ph *PersonHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id64, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		slog.DebugContext(r.Context(), "Failed parse id to number (in Delete)", "error", err)
		http.Error(w, "Failed parse id to number", http.StatusBadRequest)
		return
	}
	id := uint32(id64)

	if err := ph.Service.DeletePerson(r.Context(), id); err != nil {
		slog.DebugContext(r.Context(), "Failed delete person", "error", err)
		http.Error(w, "Failed delete person", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
