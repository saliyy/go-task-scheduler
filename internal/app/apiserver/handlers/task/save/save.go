package save

import (
	"log/slog"
	"net/http"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/entities"
	respPackage "task-scheduler/internal/lib/api/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name        string `json:"name" validate:"required,name"`
	IsCompleted string `json:"is_completed,omitempty"`
}

type Response struct {
	respPackage.Response
	Id int64 `json:"id"`
}

type TaskSaver interface {
	SaveTask(taskDTO *dto.CreateTaskDTO) (entity *entities.TaskEntity, err error)
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request")

			render.JSON(w, r, respPackage.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request")

			render.JSON(w, r, validateErr)

			return
		}

		createDto := &dto.CreateTaskDTO{
			Name:        req.Name,
			IsCompleted: false,
		}

		taskEntity, err := taskSaver.SaveTask(createDto)

		if err != nil {
			log.Error("error to save task")

			render.JSON(w, r, "error to save task")

			return
		}

		render.JSON(w, r, Response{
			Response: respPackage.OK(),
			Id:       taskEntity.Id,
		})
	}
}
