package userrepo

import (
	"fmt"
	user_dto "task-scheduler/internal/app/dto/user"
	"task-scheduler/internal/app/entities"
	"task-scheduler/internal/app/storage"
	"task-scheduler/internal/app/storage/sqlite"

	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	storage *sqlite.Storage
}

func New(sqliteStorage *sqlite.Storage) *UserRepository {
	return &UserRepository{
		storage: sqliteStorage,
	}
}

func (repo *UserRepository) CreateUser(userDto *user_dto.UserCreateDTO) (*entities.User, error) {
	const op = "storage.sqlite.userrepo.CreateUser"

	stmt, err := repo.storage.DB.Prepare("INSERT INTO users(Name, Password, Email) VALUES (?, ?, ?)")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(userDto.Name, userDto.Password, userDto.Email)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, storage.UserWithSuchEmailAlreadyExists
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entities.User{
		Id:       id,
		Name:     userDto.Name,
		Password: userDto.Password,
		Email:    userDto.Email,
	}, nil
}
