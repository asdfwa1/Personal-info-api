package service

import (
	"Test_Task_EffMob/internal/enricher"
	"Test_Task_EffMob/internal/models"
	"Test_Task_EffMob/internal/repository"
	"context"
	"log/slog"
)

type PersonService struct {
	Repo     *repository.PersonRepository
	Enricher *enricher.Enricher
}

func NewPersonService(repo *repository.PersonRepository, enricher *enricher.Enricher) *PersonService {
	return &PersonService{
		Repo:     repo,
		Enricher: enricher,
	}
}

func (ps *PersonService) CreatePerson(ctx context.Context, p *models.Person) error {
	slog.DebugContext(ctx, "Enriching person", "name", p.Name)

	enriched, err := ps.Enricher.Enrich(ctx, *p)
	if err != nil {
		slog.DebugContext(ctx, "Enrichment failed", "error", err)
		return err
	}

	slog.DebugContext(ctx, "Saving person", "person", enriched)
	return ps.Repo.Save(ctx, &enriched)
}

func (ps *PersonService) GetPersons(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, int, error) {
	slog.DebugContext(ctx, "Get person with filter", filter, "and pagination", pagination)
	return ps.Repo.FindAll(ctx, filter, pagination)
}

func (ps *PersonService) UpdatePerson(ctx context.Context, id uint32, p *models.Person) error {
	slog.DebugContext(ctx, "Update person with ID", id)
	return ps.Repo.Update(ctx, id, p)
}

func (ps *PersonService) DeletePerson(ctx context.Context, id uint32) error {
	slog.DebugContext(ctx, "Delete person with ID", id)
	return ps.Repo.Delete(ctx, id)
}
