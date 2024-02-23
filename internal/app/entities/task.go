package entities

type TaskEntity struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"date"`
	IsCompleted bool   `json:"isCompleted"`
}
