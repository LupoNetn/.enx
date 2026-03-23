package organization

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/.enx/internal/db"
	"github.com/luponetn/.enx/internal/utils"
)

type Handler struct {
	service Svc
}

func NewOrganizationHandler(service Svc) *Handler {
	return &Handler{
		service: service,
	}
}

//implement organization handler

func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
  var req CreateOrganizationRequest
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	slog.Error("could not properly pars request body,there are missing or invalid fields")
	utils.WriteError(w,http.StatusBadRequest,"invalid body sent")
	return
  }

  //create hash for passkey
  passkeyHash, err := utils.HashPassword(req.Passkey)
  if err != nil {
	slog.Error("Could not hash organization passkey", "err", err)
	return
  }

  //convert string to uuid
  parsedUUID, err := utils.StringToUUID(req.CreatedBy)

  params := db.CreateOrganizationParams{
	Name: req.Name,
	Email: req.Email,
	CreatedBy: parsedUUID,
	Passkey: passkeyHash,
  }

  //create context for db operations
  ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
  defer cancel()

  organization, err := h.service.CreateOrganization(ctx, params)
  if err != nil {
	slog.Error("Could not create a new Organization", "err", err)
	utils.WriteError(w,http.StatusInternalServerError, "could not create new organization!")
	return
  }


  utils.WriteJSON(w, http.StatusCreated,map[string]any{
	"message": "succesfully created a new organization",
     "data": organization,
  })

}

func (h *Handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	var req UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body sent")
		return
	}

	var passkeyHash *string
	if req.Passkey != nil {
		hash, err := utils.HashPassword(*req.Passkey)
		if err != nil {
			slog.Error("Could not hash organization passkey", "err", err)
			return
		}
		passkeyHash = &hash
	}

	params := db.UpdateOrganizationParams{
		ID: parsedUUID,
	}
	if req.Name != nil {
		params.Name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.Email != nil {
		params.Email = pgtype.Text{String: *req.Email, Valid: true}
	}
	if passkeyHash != nil {
		params.Passkey = pgtype.Text{String: *passkeyHash, Valid: true}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organization, err := h.service.UpdateOrganization(ctx, params)
	if err != nil {
		slog.Error("Could not update Organization", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not update organization")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully updated organization",
		"data":    organization,
	})
}

func (h *Handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.service.DeleteOrganization(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not delete Organization", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not delete organization")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully deleted organization",
	})
}

func (h *Handler) GetAllOrganizationsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	parsedUUID, err := utils.StringToUUID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	organizations, err := h.service.GetAllOrganizationsByUser(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not get organizations", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch organizations")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched organizations",
		"data":    organizations,
	})
}

func (h *Handler) GetAllUsersInOrganization(w http.ResponseWriter, r *http.Request) {
	orgID := r.PathValue("id")
	parsedUUID, err := utils.StringToUUID(orgID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid organization id")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := h.service.GetAllUsersInOrganization(ctx, parsedUUID)
	if err != nil {
		slog.Error("Could not get users in organization", "err", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch members")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "successfully fetched organization members",
		"data":    users,
	})
}
