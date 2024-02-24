package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task-scheduler/internal/app/entities"
	"task-scheduler/internal/app/storage"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"
)

type request struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type TokenResponse struct {
	Token string `json:"accessToken"`
}

type UserByEmailGiver interface {
	GetByEmail(email string) (*entities.User, error)
}

// todo move to github secrets
const SecretKey = "my.secret.key"

func Authorize(log *slog.Logger, userByEmailGiver UserByEmailGiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.auth.Authorize"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req request

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

		user, err := userByEmailGiver.GetByEmail(req.Email)

		if err != nil {

			if errors.Is(err, storage.NoUserByEmail) {

				log.Error("No user by email", req.Email)

				w.WriteHeader(http.StatusNotFound)

				render.JSON(w, r, "No user by email")

				return
			}

			fmt.Println(err.Error())

			log.Error("error to check user by email", req.Email)

			w.WriteHeader(http.StatusInternalServerError)

			render.JSON(w, r, "server error")

			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			render.JSON(w, r, "invalid password")

			return
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    strconv.Itoa(int(user.Id)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		})

		token, _ := claims.SignedString([]byte(SecretKey))

		render.JSON(w, r, TokenResponse{
			Token: token,
		})
	}
}
