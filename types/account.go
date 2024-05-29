package types

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        int       `bun:"id,pk,autoincrement"`
	UserID    uuid.UUID `bun:"type:uuid"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `bun:"default:'now()'"`
}
