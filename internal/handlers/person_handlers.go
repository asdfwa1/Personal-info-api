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
	GetPersons(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, int, error)
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

// Create
//	@Summary		Create a new person
//	@Description	Create a new person with details.
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			person	body		models.PersonInput	true	"Person Details"
//	@Success		201		{object}	models.Person		"Person created successfully"
//	@Failure		400		{object}	string				"Invalid input"
//	@Failure		500		{object}	string				"Internal server error"
//	@Router			/people/ [post]
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

// List
//	@Summary		Get List of persons with optional filters
//	@Description	Get a List of persons with optional name, surname, patronymic filters and pagination limit, offset, sorteBy, sorteOrder.
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			name		query		string					false					"Filter by name"
//	@Param			surname		query		string					false					"Filter by surname"
//	@Param			patronymic	query		string					false					"Filter by patronymic"
//	@Param			limit		query		int						false					"Limit results"			default(10)
//	@Param			offset		query		int						false					"Offset for pagination"	default(0)
//	@Param			sortBy		query		string					false					"Field to sort by (name, surname, age, etc.)"
//	@Param			sortOrder	query		string					false					"Sort order (ASC or DESC)"
//	@Success		200			{object}	map[string]interface{}	"Successful response"	SchemaExample({"data": [{"id":1,"name":"John"}], "total": 100})
//	@Failure		400			{object}	string					"Invalid query parameters"
//	@Failure		500			{object}	string					"Internal server error"
//	@Router			/people/ [get]
func (ph *PersonHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterParams{
		Name:       r.URL.Query().Get("name"),
		Surname:    r.URL.Query().Get("surname"),
		Patronymic: r.URL.Query().Get("patronymic"),
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		} else {
			slog.DebugContext(r.Context(), "Invalid limit parameter", "error", err)
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		} else {
			slog.DebugContext(r.Context(), "Invalid offset parameter", "error", err)
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	sortBy := r.URL.Query().Get("sortBy")
	sortOrder := r.URL.Query().Get("sortOrder")

	validSortFields := map[string]bool{
		"id": true, "name": true, "surname": true,
		"patronymic": true, "age": true, "gender": true, "nationality": true,
	}
	if !validSortFields[sortBy] && sortBy != "" {
		slog.DebugContext(r.Context(), "Invalid sortBy parameter", "sortBy", sortBy)
		http.Error(w, "Invalid sort field", http.StatusBadRequest)
		return
	}
	if sortOrder != "ASC" && sortOrder != "DESC" && sortOrder != "" {
		slog.DebugContext(r.Context(), "Invalid sortOrder parameter", "sortOrder", sortOrder)
		http.Error(w, "sortOrder must be ASC or DESC", http.StatusBadRequest)
		return
	}

	pagination := models.PaginationParams{
		Limit:     limit,
		OffSet:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	persons, total, err := ph.Service.GetPersons(r.Context(), filter, pagination)
	if err != nil {
		slog.DebugContext(r.Context(), "Failed getting list person", "error", err)
		http.Error(w, "Failed getting list person", http.StatusInternalServerError)
		return
	}

	response := struct {
		Data  []models.Person `json:"data"`
		Total int             `json:"total"`
	}{
		Data:  persons,
		Total: total,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

//	Update
//
//	@Summary		Update an existing person
//	@Description	Update the details of an existing person by ID.
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Person ID"
//	@Param			person	body		models.Person	true	"Updated Person Details"
//	@Success		200		{string}	string			"Update person with ID"
//	@Failure		400		{object}	string			"Invalid input or ID"
//	@Failure		404		{object}	string			"Person not found"
//	@Failure		500		{object}	string			"Internal server error"
//	@Router			/people/{id} [put]
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

// Delete
//	@Summary		Delete a person by ID
//	@Description	Delete the person with the given ID from the system.
//	@Tags			persons
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Person ID"
//	@Success		204	{string}	string	"Person deleted successfully"
//	@Failure		400	{object}	string	"Invalid ID"
//	@Failure		404	{object}	string	"Person not found"
//	@Failure		500	{object}	string	"Internal server error"
//	@Router			/people/{id} [delete]
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
