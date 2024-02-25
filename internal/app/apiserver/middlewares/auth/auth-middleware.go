package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"task-scheduler/internal/app/apiserver/handlers/auth/oauth"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
)

const UserIdContextKey = "userId"

func CurrentUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headers := r.Header

		_, has := headers["Authorization"]

		if !has {
			render.JSON(w, r, "No Authorization header")

			return
		}

		bearer := r.Header.Get("Authorization")

		if bearer == "" {
			render.JSON(w, r, "No auth token in Authorization header")

			return
		}

		parts := strings.Split(bearer, " ")

		if len(parts) < 2 {
			render.JSON(w, r, "No auth token in Authorization header")

			return
		}

		accessToken := parts[1]

		token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(oauth.SecretKey), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			render.JSON(w, r, "Unauthorized")

			return
		}

		claims := token.Claims.(*jwt.StandardClaims)

		userId, _ := strconv.Atoi(claims.Issuer)

		ctx := context.WithValue(r.Context(), UserIdContextKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
