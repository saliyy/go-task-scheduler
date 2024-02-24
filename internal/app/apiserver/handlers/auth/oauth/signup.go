package oauth

import (
	"errors"
	"net/http"
	user_dto "task-scheduler/internal/app/dto/user"
	"task-scheduler/internal/app/entities"
	"task-scheduler/internal/app/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type Request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Response struct {
	status int
	Id     int64 `json:"id"`
}

type UserCreator interface {
	CreateUser(createDto *user_dto.UserCreateDTO) (*entities.User, error)
}

func New(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.auth.signup.New"

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

			render.JSON(w, r, "invalid data")

			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 15)

		createUserDto := &user_dto.UserCreateDTO{
			Name:     req.Name,
			Password: string(password),
			Email:    req.Email,
		}

		user, err := userCreator.CreateUser(createUserDto)

		if err != nil {
			if errors.Is(err, storage.UserWithSuchEmailAlreadyExists) {
				log.Error("Trying to create user with email that already has in system")

				w.WriteHeader(http.StatusUnprocessableEntity)

				render.JSON(w, r, "User with such email already exists!")

				return
			}

			log.Error("Error to create user")

			w.WriteHeader(http.StatusUnprocessableEntity)

			render.JSON(w, r, "Error to create user")

			return

		}

		log.Info("user created", user)

		render.JSON(w, r, Response{
			status: http.StatusCreated,
			Id:     user.Id,
		})
	}
}
