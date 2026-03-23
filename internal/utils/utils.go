// internal/utils/uuid.go
package utils

import (
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

// converts a string to pgtype.UUID for database operations
func StringToUUID(s string) (pgtype.UUID, error) {
    u, err := uuid.Parse(s)
    if err != nil {
        return pgtype.UUID{}, fmt.Errorf("invalid uuid: %w", err)
    }

    var id pgtype.UUID
    copy(id.Bytes[:], u[:])
    id.Valid = true
    return id, nil
}

// converts pgtype.UUID back to a plain string
func UUIDToString(id pgtype.UUID) string {
    u, err := uuid.FromBytes(id.Bytes[:])
    if err != nil {
        return ""
    }
    return u.String()
}