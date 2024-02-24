package storage

import "errors"

var (
	TaskAlreadyExists              = errors.New("TaskAlreadyExists")
	NotSuchTask                    = errors.New("NotSuchTask")
	UserWithSuchEmailAlreadyExists = errors.New("UserWithSuchEmailAlreadyExists")
	NoUserByEmail                  = errors.New("NoUserByEmail")
)
