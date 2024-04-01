package listrepo

import (
	"context"
	"fmt"
	dto "task-scheduler/internal/app/dto/list"
	"task-scheduler/internal/app/entities"
	"task-scheduler/internal/app/storage/sqlite"
)

type ListRepository struct {
	storage *sqlite.Storage
}

func New(sqliteStorage *sqlite.Storage) *ListRepository {
	return &ListRepository{
		storage: sqliteStorage,
	}
}

func (r *ListRepository) Сreate(dto *dto.ListCreateDTO) (entity *entities.ListEntity, err error) {
	const op = "storage.sqlite.list.create"

	context := context.TODO()
	tx, err := r.storage.DB.BeginTx(context, nil)

	defer tx.Rollback()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO lists(Name) VALUES (?)")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(dto.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = tx.Prepare("INSERT INTO users_lists(UserId, ListId) VALUES (?, ?)")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(dto.UserId, id)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entities.ListEntity{
		Id:   id,
		Name: dto.Name,
	}, nil
}
