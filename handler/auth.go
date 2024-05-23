package handler

import (
	"dreampicai/view/auth"
	"net/http"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}
