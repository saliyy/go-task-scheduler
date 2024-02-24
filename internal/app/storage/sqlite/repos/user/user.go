package userrepo

import (
	"database/sql"
	"errors"
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

func (repo *UserRepository) GetByEmail(email string) (*entities.User, error) {
	const op = "storage.sqlite.userrepo.GetUserByEmail"

	stmt, err := repo.storage.DB.Prepare("SELECT * FROM users WHERE Email = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	var userEntity entities.User

	err = stmt.QueryRow(email).Scan(&userEntity.Id, &userEntity.Name, &userEntity.Password, &userEntity.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.NoUserByEmail
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return &userEntity, nil
}
