package list

import (
	"encoding/json"
	"net/http"
	"task-scheduler/internal/app/apiserver/middlewares/auth"
	"task-scheduler/internal/app/entities"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TaskGiver interface {
	GetTasksByUserId(userId int) (entity []entities.TaskEntity, err error)
}

func New(log *slog.Logger, taskGiver TaskGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.list.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		currentUserId := r.Context().Value(auth.UserIdContextKey).(int)

		list, err := taskGiver.GetTasksByUserId(currentUserId)

		if err != nil {

			log.Error("error to save task")

			render.JSON(w, r, "error to load tasks")

			return
		}

		marshalledList, err := json.Marshal(list)
		if err != nil {
			log.Error("error to unmarshal tasks")

			render.JSON(w, r, "error to load tasks")

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(marshalledList)
	}
}
