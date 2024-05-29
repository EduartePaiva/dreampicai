package types

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        int       `bun:"id,pk,autoincrement"`
	UserID    uuid.UUID `bun:"user_id,type:uuid"`
	UserName  string    `bun:"user_name"`
	CreatedAt time.Time `bun:"created_at,default:'now()'"`
}
