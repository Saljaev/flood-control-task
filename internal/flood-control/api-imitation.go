package flood_control

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type Error struct {
	ErrorMessage string `json:"error_message"`
}

type Response struct {
	Message string `json:"message"`
}

func WriteFailure(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)

	data, _ := json.Marshal(Error{ErrorMessage: msg})
	w.Write(data)
}

func WriteSuccess(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusOK)

	data, _ := json.Marshal(Response{Message: msg})
	w.Write(data)
}

func (h *Handler) Nothing(ctx context.Context, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "FloodHandler - Nothing"

		log := log.With(
			slog.String("op", op),
		)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		ID := r.Header.Get("user_id")
		if ID == "" {
			log.Error("invalid userID", slog.Any("ID", ID))
			WriteFailure(w, http.StatusBadRequest, "login to your account")
			return
		}

		userID, err := strconv.Atoi(ID)
		if err != nil || userID <= 0 {
			log.Error("invalid userID", slog.Any("ID", ID))
			WriteFailure(w, http.StatusBadRequest, "login to your account")
			return
		}

		check, err := h.FloodControl.Check(ctx, int64(userID))
		if err != nil {
			log.Error("failed to check user rate", slog.Any("error", err))
			WriteFailure(w, http.StatusInternalServerError, "server error")
			return
		}

		if !check {
			log.Info("user has max available requests")
			WriteFailure(w, http.StatusForbidden, "too many requests, try again later")
			return
		} else {
			log.Info("request from user", slog.Any("ID", userID))
			WriteSuccess(w, "enjoy our API")
			return
		}
	}
}
