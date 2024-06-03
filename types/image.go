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
	ID            uuid.UUID `bun:"id,pk,default:gen_random_uuid()"`
	UserID        uuid.UUID `bun:"user_id,type:uuid"`
	ImageLocation string
	Status        ImageStatus
	BatchID       uuid.UUID
	Prompt        string
	Deleted       bool      `bun:"default:'false'"`
	CreatedAt     time.Time `bun:"default:'now()'"`
	DeletedAt     time.Time `bun:",nullzero"`
}
