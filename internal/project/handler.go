package project

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

func NewProjectHandler(service Svc) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("could not properly pars request body,there are missing or invalid fields")
		utils.WriteError(w, http.StatusBadRequest, "invalid body sent")
		return
	}

	passkeyHash, err := utils.HashPassword(req.Passkey)
	if err != nil {
		slog.Error("Could not hash project passkey", "err", err)
		return
	}

	parsedOrgUUID, err := utils.StringToUUID(req.OrganizationID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	parsedCreatedByUUID, err := utils.StringToUUID(req.CreatedBy)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid created_by id")
		return
	}

	params := db.CreateProjectParams{
		Name:           req.Name,
		Passkey:        passkeyHash,
		OrganizationID: parsedOrgUUID,
		CreatedBy:      parsedCreatedByUUID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//check if project name already exists in the organization
	_, err = h.service.GetProjectByName(ctx, db.GetProjectByNameParams{
		Name:           req.Name,
		OrganizationID: parsedOrgUUID,
	})
	if err == nil {
		utils.WriteError(w, http.StatusConflict, "project name already exists in this organization")
		return
	}

	project, err := h.service.CreateProject(ctx, params)
	if err != nil {
		slog.Error("Could not create a new Project", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create new project!")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "succesfully created a new project",
		"data":    project,
	})
}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	var req UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body sent")
		return
	}

	var passkeyHash *string
	if req.Passkey != nil {
		hash, err := utils.HashPassword(*req.Passkey)
		if err != nil {
			slog.Error("Could not hash project passkey", "err", err)
			return
		}
		passkeyHash = &hash
	}

	params := db.UpdateProjectParams{
		ID: parsedUUID,
	}
	if req.Name != nil {
		params.Name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if passkeyHash != nil {
		params.Passkey = pgtype.Text{String: *passkeyHash, Valid: true}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	project, err := h.service.UpdateProject(ctx, params)
	if err != nil {
		slog.Error("Could not update Project", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not update project")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully updated project",
		"data":    project,
	})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.service.DeleteProject(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not delete Project", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not delete project")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully deleted project",
	})
}

func (h *Handler) GetProjectsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	parsedUUID, err := utils.StringToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projects, err := h.service.GetProjectsByUser(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not get projects", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch projects")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched projects",
		"data":    projects,
	})
}

func (h *Handler) GetAllUsersInProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := h.service.GetAllUsersInProject(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not get users in project", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch members")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched project members",
		"data":    users,
	})
}

func (h *Handler) GetProjectByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	orgID := r.URL.Query().Get("organization_id")

	if name == "" {
		utils.WriteError(w, http.StatusBadRequest, "project name is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var project interface{}
	var err error

	if orgID != "" {
		parsedOrgUUID, errUUID := utils.StringToUUID(orgID)
		if errUUID != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
			return
		}
		project, err = h.service.GetProjectByName(ctx, db.GetProjectByNameParams{
			Name:           name,
			OrganizationID: parsedOrgUUID,
		})
	} else {
		userID, errAuth := auth.GetUserIDFromContext(r.Context())
		if errAuth != nil {
			utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		project, err = h.service.GetProjectByNameForUser(ctx, db.GetProjectByNameForUserParams{
			Name:   name,
			UserID: userID,
		})
	}

	if err != nil {
		slog.Error("Could not get project by name", "err", err)
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched project",
		"data":    project,
	})
}
