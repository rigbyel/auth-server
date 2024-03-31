package login

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/auth-server/internal/lib/jwt"
	"github.com/rigbyel/auth-server/internal/lib/response"
	"github.com/rigbyel/auth-server/internal/models"
	"github.com/rigbyel/auth-server/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	AccessToken string `json:"access_token"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProvider interface {
	User(login string) (*models.User, error)
}

// New create a HandlerFunc to handle /login endpoint
func New(log *slog.Logger, userProvider UserProvider, authSecret string, tokenTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.Login.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		// decoding request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			response.RespondWithError(w, 400, "invalid request")
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// getting user info from storage
		user, err := userProvider.User(req.Email)
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found", slog.String("user", req.Email))

			response.RespondWithError(w, 404, "user not found")

			return
		}
		if err != nil {
			log.Error("error finding user", slog.String("error", err.Error()))

			response.RespondWithError(w, 500, "internal error")

			return
		}

		// checking user's password
		if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(req.Password)); err != nil {
			log.Info("invalid credentials", slog.String("error", err.Error()))

			response.RespondWithError(w, 401, "invalid credentials")

			return
		}

		log.Info("user logged in successfully")

		// creating new jwt token
		token, err := jwt.NewToken(*user, authSecret, tokenTL)
		if err != nil {
			log.Error("failed to generate token", slog.String("string", err.Error()))

			response.RespondWithError(w, 500, "internal error")

			return
		}

		response.RespondWithJson(w, 200,
			Response{
				AccessToken: token,
			},
		)
	}
}
