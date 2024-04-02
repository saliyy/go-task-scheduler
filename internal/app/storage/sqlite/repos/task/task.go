package taskrepo

import (
	"fmt"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/entities"
	"task-scheduler/internal/app/storage/sqlite"
	"time"
)

type TaskRepository struct {
	storage *sqlite.Storage
}

func New(sqliteStorage *sqlite.Storage) *TaskRepository {
	return &TaskRepository{
		storage: sqliteStorage,
	}
}

func (r *TaskRepository) SaveTask(taskDTO *dto.CreateTaskDTO) (entity *entities.TaskEntity, err error) {
	const op = "storage.sqlite.saveTask"

	stmt, err := r.storage.DB.Prepare("INSERT INTO tasks(Name, UserId, ListId) VALUES (?, ?, ?)")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(taskDTO.Name, taskDTO.UserId, taskDTO.ListId)
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
		UserId:      taskDTO.UserId,
		ListId:      taskDTO.ListId,
	}, nil
}

func (r *TaskRepository) GetTasksByUserId(userId int) (entites []entities.TaskEntity, err error) {
	const op = "storage.sqlite.GetTasks"

	rows, err := r.storage.DB.Query("SELECT * FROM tasks WHERE UserId = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var tasks []entities.TaskEntity

	for rows.Next() {
		var task entities.TaskEntity
		if err := rows.Scan(&task.Id, &task.Name, &task.CreatedAt, &task.IsCompleted, &task.UserId); err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil

}
