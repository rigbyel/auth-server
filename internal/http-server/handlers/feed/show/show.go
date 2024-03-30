package show

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rigbyel/auth-server/internal/lib/jwt"
	"github.com/rigbyel/auth-server/internal/lib/response"
)

func New(log *slog.Logger, authSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.feed.show.New"

		// setting up logger
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// check if there's a valid token in the request
		tokenString := r.Header.Get("Authorization-access")
		_, err := jwt.GetTokenClaims(tokenString, authSecret)
		if err != nil {
			response.RespondWithError(w, 401, "unauthorized")

			return
		}

		response.RespondWithJson(w, 200, nil)
	}
}
