package list

import (
	"fmt"
	"net/http"
	"task-scheduler/internal/app/apiserver/middlewares/auth"
	dto "task-scheduler/internal/app/dto/list"

	"golang.org/x/exp/slog"

	"task-scheduler/internal/app/entities"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	status int
	Id     int64 `json:"id"`
}

type ListCreator interface {
	Сreate(dto *dto.ListCreateDTO) (entity *entities.ListEntity, err error)
}

func New(log *slog.Logger, listCreator ListCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.list.save.New"

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

		createDto := &dto.ListCreateDTO{
			Name:   req.Name,
			UserId: currentUserId,
		}

		listEntity, err := listCreator.Сreate(createDto)

		fmt.Println(err)

		if err != nil {
			log.Error("error to save list")

			render.JSON(w, r, "error to save list")

			return
		}

		w.WriteHeader(http.StatusCreated)

		render.JSON(w, r, Response{
			status: http.StatusCreated,
			Id:     listEntity.Id,
		})
	}
}
