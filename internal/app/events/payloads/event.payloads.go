package payloads

import (
	"task-scheduler/internal/app/entities"
)

type UserCreatedPayload struct {
	UserEntity entities.User
}
