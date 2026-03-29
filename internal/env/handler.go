package env

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/enx/internal/auth"
	"github.com/luponetn/enx/internal/db"
	"github.com/luponetn/enx/internal/utils"
)

type Handler struct {
	service Svc
}

func NewHandler(service Svc) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateEnv(w http.ResponseWriter, r *http.Request) {
	var req CreateEnvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Could not parse create env body request", "err", err)
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	parsedProjectUUID, err := utils.StringToUUID(req.ProjectID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Double check if name exists in project
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = h.service.GetEnvByNameInProject(ctx, db.GetEnvByNameInProjectParams{
		Name:      req.Name,
		ProjectID: parsedProjectUUID,
	})
	if err == nil {
		utils.WriteError(w, http.StatusConflict, "an environment with this name already exists in the project")
		return
	}

	params := db.CreateEnvParams{
		Name:        req.Name,
		ProjectID:   parsedProjectUUID,
		Variables:   []byte(req.Variables),
		Description: req.Description,
		CreatedBy:   userID,
	}

	env, err := h.service.CreateEnv(ctx, params)
	if err != nil {
		slog.Error("Could not create environment", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to create environment")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "successfully created environment",
		"data":    env,
	})
}

func (h *Handler) GetEnvsByProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("project_id")
	parsedUUID, err := utils.StringToUUID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	envs, err := h.service.GetEnvsByProject(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not fetch environments", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch environments")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched environments",
		"data":    envs,
	})
}

func (h *Handler) GetEnvByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid environment id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	env, err := h.service.GetEnvByID(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not fetch environment", "err", err)
		utils.WriteError(w, http.StatusNotFound, "environment not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched environment",
		"data":    env,
	})
}

func (h *Handler) UpdateEnv(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid environment id")
		return
	}

	var req UpdateEnvRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	params := db.UpdateEnvParams{
		ID: parsedUUID,
	}
	if req.Name != nil {
		params.Name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.Description != nil {
		params.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.Variables != nil {
		params.Variables = []byte(*req.Variables)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	env, err := h.service.UpdateEnv(ctx, params)
	if err != nil {
		slog.Error("Could not update environment", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to update environment")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully updated environment",
		"data":    env,
	})
}

func (h *Handler) DeleteEnv(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid environment id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.service.DeleteEnv(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not delete environment", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "failed to delete environment")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully deleted environment",
	})
}
