package enricher

import (
	"Test_Task_EffMob/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"
)

type Enricher struct {
	Client *http.Client
}

func NewEnricher() *Enricher {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &Enricher{
		Client: client,
	}
}

func (e *Enricher) Enrich(ctx context.Context, p models.Person) (models.Person, error) {
	name := p.Name
	type ResponseAge struct {
		Age uint8 `json:"age"`
	}

	type ResponseGender struct {
		Gender      string  `json:"gender"`
		Probability float32 `json:"probability"`
	}

	type ResponseNational struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float32 `json:"probability"`
		}
	}

	agifyURL := os.Getenv("Agify_URL")
	ageURLwithName := fmt.Sprintf(agifyURL+"%s", name)
	ageRes, err := e.Client.Get(ageURLwithName)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			slog.WarnContext(ctx, "Timeout calling Agify API")
			return p, nil
		} else {
			return p, fmt.Errorf("error calling Agify API: %w", err)
		}
	}
	defer func() {
		if err := ageRes.Body.Close(); err != nil {
			slog.Error("Failed to close ageBody", "error", err)
		}
	}()
	var age ResponseAge
	if err := json.NewDecoder(ageRes.Body).Decode(&age); err == nil {
		p.Age = age.Age
	}

	genderizeURL := os.Getenv("Genderize_URL")
	genderURLwithName := fmt.Sprintf(genderizeURL+"%s", name)
	genderRes, err := e.Client.Get(genderURLwithName)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			slog.WarnContext(ctx, "Timeout calling Genderize API")
			return p, nil
		} else {
			return p, fmt.Errorf("error calling Genderize API: %w", err)
		}
	}
	defer func() {
		if err := genderRes.Body.Close(); err != nil {
			slog.Error("Failed to close genderBody", "error", err)
		}
	}()
	var gender ResponseGender
	if err := json.NewDecoder(genderRes.Body).Decode(&gender); err == nil {
		p.Gender = gender.Gender
	}

	nationalizeURL := os.Getenv("Nationalize_URL")
	nationalityURLwithName := fmt.Sprintf(nationalizeURL+"%s", name)
	nationalityRes, err := e.Client.Get(nationalityURLwithName)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			slog.WarnContext(ctx, "Timeout calling Nationalize API")
			return p, nil
		} else {
			return p, fmt.Errorf("error calling Nationalize API: %w", err)
		}
	}
	defer func() {
		if err := nationalityRes.Body.Close(); err != nil {
			slog.Error("Failed to close nationalityBody", "error", err)
		}
	}()
	var nationality ResponseNational
	if err := json.NewDecoder(nationalityRes.Body).Decode(&nationality); err == nil && len(nationality.Country) > 0 {
		p.Nationality = nationality.Country[0].CountryID
	}

	slog.InfoContext(ctx, "Person enriched", "age:", p.Age, "gender", p.Gender, "nationality:", p.Nationality, "Name:", name)
	return p, nil
}
