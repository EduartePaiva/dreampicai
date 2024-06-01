package types

import (
	"time"

	"github.com/google/uuid"
)

type ImageStatus int

const (
	ImageStatusFailed ImageStatus = iota
	ImageStatusPending
	ImageStatusCompleted
)

type Image struct {
	ID        int `bun:"id,pk,autoincrement"`
	UserID    uuid.UUID
	Username  string
	Status    ImageStatus
	CreatedAt time.Time `bun:"default:'now()'"`
}
