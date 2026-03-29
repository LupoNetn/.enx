// internal/utils/uuid.go
package utils

import (
    "github.com/google/uuid"
)

// converts a string to uuid.UUID for database operations
func StringToUUID(s string) (uuid.UUID, error) {
    return uuid.Parse(s)
}

// converts uuid.UUID back to a plain string
func UUIDToString(id uuid.UUID) string {
    return id.String()
}
