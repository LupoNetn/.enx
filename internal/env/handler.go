package env

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/luponetn/enx/internal/utils"
)

type Handler struct {
	service Svc
}

func NewHandler(service Svc) *Handler {
	return &Handler{service: service}
}

// implement handlers for creating of env variables
func (h *Handler) CreateEnv(w http.ResponseWriter, r *http.Request) {
	var req CreateEnvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Could not parse create env body request")
		utils.WriteError(w,http.StatusBadRequest,"invalid request body")
		return
	}
}
