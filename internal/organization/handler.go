package organization

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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
