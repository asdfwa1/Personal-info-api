package repository

import (
	"Test_Task_EffMob/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
)

var ErrNotFound = errors.New("not found person with id")

type PersonRepository struct {
	SaveStmt   *sql.Stmt
	UpdateStmt *sql.Stmt
	DeleteStmt *sql.Stmt
	FindStmt   string
	DB         *sql.DB
}

func NewPersonRepository(db *sql.DB) (*PersonRepository, error) {
	saveStmt, err := db.Prepare(`INSERT INTO persons (name, surname, patronymic, age, gender, nationality)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`)
	if err != nil {
		return nil, err
	}

	updateStmt, err := db.Prepare(`UPDATE persons
		SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
		WHERE id = $7;`)
	if err != nil {
		return nil, err
	}

	deleteStmt, err := db.Prepare(`DELETE FROM persons WHERE id = $1;`)
	if err != nil {
		return nil, err
	}

	return &PersonRepository{
		SaveStmt:   saveStmt,
		UpdateStmt: updateStmt,
		DeleteStmt: deleteStmt,
		FindStmt:   `SELECT id, name, surname, patronymic, age, gender, nationality FROM persons WHERE 1=1`,
		DB:         db,
	}, nil
}

func (pr *PersonRepository) Save(ctx context.Context, p *models.Person) error {
	return pr.SaveStmt.QueryRowContext(ctx, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality).Scan(&p.ID)
}

func (pr *PersonRepository) Update(ctx context.Context, id uint32, p *models.Person) error {
	ex, err := pr.UpdateStmt.ExecContext(ctx, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Nationality, id)
	if err != nil {
		return err
	}
	rowsAffected, err := ex.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		slog.DebugContext(ctx, "No updated rows with", "id", id)
		return ErrNotFound
	}
	return nil
}

func (pr *PersonRepository) Delete(ctx context.Context, id uint32) error {
	ex, err := pr.DeleteStmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	rowsAffected, err := ex.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		slog.DebugContext(ctx, "No deleted rows with", "id", id)
		return ErrNotFound
	}
	return nil
}

func (pr *PersonRepository) FindAll(ctx context.Context, filter models.FilterParams, pagination models.PaginationParams) ([]models.Person, error) {
	query := pr.FindStmt
	args := []any{}
	i := 1

	if filter.Name != "" {
		query = query + fmt.Sprintf(" AND name = $%d", i)
		args = append(args, filter.Name)
		i++
	}
	if filter.Surname != "" {
		query = query + fmt.Sprintf(" AND surname = $%d", i)
		args = append(args, filter.Surname)
		i++
	}
	if filter.Patronymic != "" {
		query = query + fmt.Sprintf(" AND patronymic = $%d", i)
		args = append(args, filter.Patronymic)
		i++
	}

	if pagination.Limit != 0 {
		query = query + " LIMIT $" + strconv.Itoa(i)
		args = append(args, pagination.Limit)
		i++
	}

	if pagination.OffSet != 0 {
		query = query + " OFFSET $" + strconv.Itoa(i)
		args = append(args, pagination.OffSet)
		i++
	}

	rows, err := pr.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var p models.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Nationality); err != nil {
			return nil, err
		}
		people = append(people, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

func (pr *PersonRepository) Close() error {
	if err := pr.SaveStmt.Close(); err != nil {
		return err
	}
	if err := pr.UpdateStmt.Close(); err != nil {
		return err
	}
	if err := pr.DeleteStmt.Close(); err != nil {
		return err
	}
	return nil
}
