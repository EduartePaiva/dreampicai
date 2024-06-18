package handler

import (
	"dreampicai/view/home"
	"fmt"
	"net/http"
	"time"
)

func HandleLongProcess(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(time.Second * 5)
	return render(r, w, home.UserLikes(1000))
}

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	fmt.Printf("%+v\n", user.Account)
	return home.Index().Render(r.Context(), w)
}
