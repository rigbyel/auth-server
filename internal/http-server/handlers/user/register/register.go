package register

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rigbyel/auth-server/internal/lib/response"
	validate "github.com/rigbyel/auth-server/internal/lib/validate/password"
	"github.com/rigbyel/auth-server/internal/models"
	"github.com/rigbyel/auth-server/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type Response struct {
	ID int64 `json:"id"`
}

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSaver interface {
	SaveUser(u *models.User) (*models.User, error)
}

// New creates a new HandlerFunc to handle user registration
func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

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

		// validating login and password
		var validationErrs []string

		validationErrs = validate.ValidatePassword(req.Password)

		if len(validationErrs) != 0 {
			log.Error("invalid request data")

			response.RespondWithError(w, 400, strings.Join(validationErrs, ", "))
			return
		}

		// hashing password
		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", slog.String("error", err.Error()))

			response.RespondWithError(w, 500, "internal error")

			return
		}

		// creating user
		user := &models.User{
			Email:    req.Email,
			PassHash: passHash,
		}

		// saving user in the storage
		user, err = userSaver.SaveUser(user)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("user", req.Email))

			response.RespondWithError(w, 400, "user already exists")

			return
		}
		if err != nil {
			log.Error("error creating user", slog.String("error", err.Error()))

			response.RespondWithError(w, 500, "internal error")

			return
		}

		log.Info("user added", slog.Int64("id", user.Id))

		response.RespondWithJson(w, 201,
			Response{
				ID: user.Id,
			},
		)
	}
}
