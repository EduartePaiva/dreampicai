package handler

import (
	"dreampicai/db"
	"dreampicai/pkg/kit/validate"
	"dreampicai/pkg/sb"
	"dreampicai/pkg/util"
	"dreampicai/types"
	"dreampicai/view/auth"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionUserKey        = "user"
	sessionAccessTokenKey = "accessToken"
)

func HandleResetPasswordIndex(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	fmt.Println(accessToken)
	return render(r, w, auth.ResetPassword())
}
func HandleResetPasswordCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	// port := os.Getenv("HTTP_LISTEN_ADDR")
	// err := sb.MyResetPasswordForEmail(r.Context(), user.Email, "http://localhost"+port+"/auth/reset-password")
	// if err != nil {
	// 	return err
	// }
	// to avoid running out of emails

	return render(r, w, auth.ResetPasswordInitiated(user.Email))
}
func HandleResetPasswordUpdate(w http.ResponseWriter, r *http.Request) error {
	params := auth.ResetPasswordParams{
		NewPassword:     r.FormValue("new-password"),
		ConfirmPassword: r.FormValue("confirm-password"),
	}
	errors := auth.ResetPasswordErrors{}
	ok := validate.New(&params, validate.Fields{
		"NewPassword": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.NewPassword),
			validate.Message("passwords do not match"),
		),
	}).Validate(&errors)
	if !ok {
		return render(r, w, auth.ResetPasswordForm(errors))
	}

	//do logic to change psw

	// if err := render(r, w, components.Toast("Password updated")); err != nil {
	// 	return err
	// }
	return hxRedirect(w, r, "/")
}

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.AccountSetup())
}
func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.AccountSetupData{
		Username: r.FormValue("username"),
	}
	var errors auth.AccountSetupErrors
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(50)),
	}).Validate(&errors)
	if !ok {
		return render(r, w, auth.AccountSetupForm(params, errors))
	}

	user := getAuthenticatedUser(r)
	account := types.Account{
		UserID:   user.ID,
		UserName: params.Username,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")
}

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignUp())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	if ok, signupErros := util.IsValidSignupForm(params); !ok {
		return render(r, w, auth.SignupForm(params, signupErros))
	}

	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}

	return render(r, w, auth.SignupSuccess(user.Email))
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	//todo read this from the env variable
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3131/auth/callback",
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))
	}

	path := r.Header.Get("Hx-Current-Url")
	url, err := url.Parse(path)
	if err != nil {
		return hxRedirect(w, r, "/")
	}
	toRedirect := url.Query().Get("to")
	if len(toRedirect) == 0 {
		toRedirect = "/"
	}

	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		return err
	}
	return hxRedirect(w, r, toRedirect)
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = ""
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, sessionUserKey)
	session.Values[sessionAccessTokenKey] = accessToken
	return session.Save(r, w)
}
