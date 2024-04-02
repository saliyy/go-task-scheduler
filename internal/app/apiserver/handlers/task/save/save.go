package save

import (
	"net/http"
	"task-scheduler/internal/app/apiserver/middlewares/auth"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/entities"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name   string `json:"name"`
	ListId int    `json:"listId"`
}

type Response struct {
	status int
	Id     int64 `json:"id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=TaskSaver
type TaskSaver interface {
	SaveTask(taskDTO *dto.CreateTaskDTO) (entity *entities.TaskEntity, err error)
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request")

			render.JSON(w, r, "faild")

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request")

			render.JSON(w, r, "faild")

			return
		}

		currentUserId := r.Context().Value(auth.UserIdContextKey).(int)

		createDto := &dto.CreateTaskDTO{
			Name:        req.Name,
			IsCompleted: false,
			UserId:      currentUserId,
			ListId:      req.ListId,
		}

		taskEntity, err := taskSaver.SaveTask(createDto)

		if err != nil {
			log.Error("error to save task")

			render.JSON(w, r, "error to save task")

			return
		}

		w.WriteHeader(http.StatusCreated)

		render.JSON(w, r, Response{
			status: http.StatusCreated,
			Id:     taskEntity.Id,
		})
	}
}
