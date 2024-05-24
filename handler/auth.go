package handler

import (
	"dreampicai/pkg/sb"
	"dreampicai/view/auth"
	"fmt"
	"net/http"

	"github.com/nedpals/supabase-go"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	resp, err := sb.Client.Auth.SignIn()
	return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
		InvalidCredentials: "The credentials you have entered are invalid",
	}))
	fmt.Println(credentials)
	return nil
}
