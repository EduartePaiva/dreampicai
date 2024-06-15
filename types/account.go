package types

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int       `bun:"id,pk,autoincrement"`
	UserID    uuid.UUID `bun:"user_id,type:uuid"`
	UserName  string    `bun:"user_name"`
	Credits   int
	CreatedAt time.Time `bun:"created_at,default:'now()'"`
}
