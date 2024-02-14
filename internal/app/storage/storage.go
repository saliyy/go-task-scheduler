package storage

import "errors"

var (
	TaskAlreadyExists = errors.New("TaskAlreadyExists")
	NotSuchTask       = errors.New("NotSuchTask")
)
