package types

import "github.com/google/uuid"

type AuthenticatedUser struct {
	ID       uuid.UUID
	Email    string
	LoggedIn bool
	Account
}

type contextKey string

const UserContextKey contextKey = "user"
