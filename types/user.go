package types

type AuthenticatedUser struct {
	Email    string
	LoggedIn bool
}

type contextKey string

const UserContextKey contextKey = "user"
