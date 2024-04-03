package events

import (
	"task-scheduler/internal/app/listeners"

	"github.com/nuttech/bell/v2"
)

const (
	UserCreated = "UserCreated"
)

func Init(listeners *listeners.Listeners) {
	bell.Listen(UserCreated, listeners.CreateDefaultListListener.Handle())
}
