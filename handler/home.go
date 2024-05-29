package handler

import (
	"dreampicai/view/home"
	"fmt"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	fmt.Printf("%+v\n", user.Account)
	return home.Index().Render(r.Context(), w)
}
