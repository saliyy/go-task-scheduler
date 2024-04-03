package listeners

import (
	"task-scheduler/internal/app/entities"

	"github.com/nuttech/bell/v2"
)

type DefaultListCreator interface {
	CreateDefaultList(userId int)
}

type CreateDefaultList struct {
	ListCreator DefaultListCreator
}

func (l *CreateDefaultList) Handle() func(message bell.Message) {
	return func(message bell.Message) {
		user := message.(*entities.User)
		l.ListCreator.CreateDefaultList(int(user.Id))
	}
}

func NewDefaultListCreator(listCreator DefaultListCreator) *CreateDefaultList {
	return &CreateDefaultList{
		ListCreator: listCreator,
	}
}
