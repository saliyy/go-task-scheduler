package sqlite

import (
	"database/sql"
	"fmt"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/entities"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "sqlite.storage.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// todo move to migration
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS tasks(
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Name TEXT NOT NULL,
			CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			IsCompleted INTEGER NOT NULl DEFAULT 0
		)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveTask(taskDTO *dto.CreateTaskDTO) (entity *entities.TaskEntity, err error) {
	const op = "storage.sqlite.saveTask"

	stmt, err := s.db.Prepare("INSERT INTO tasks(Name) VALUES (?)")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	res, err := stmt.Exec(taskDTO.Name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entities.TaskEntity{
		Id:          id,
		Name:        taskDTO.Name,
		IsCompleted: false,
		CreatedAt:   time.DateTime,
	}, nil
}

func (s *Storage) GetTasks() (entites []entities.TaskEntity, err error) {
	const op = "storage.sqlite.GetTasks"

	rows, err := s.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var tasks []entities.TaskEntity

	for rows.Next() {
		var task entities.TaskEntity
		if err := rows.Scan(&task.Id, &task.Name, &task.CreatedAt, &task.IsCompleted); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil

}
