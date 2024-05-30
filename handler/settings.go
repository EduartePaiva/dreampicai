package handler

import (
	"dreampicai/db"
	"dreampicai/pkg/kit/validate"
	"dreampicai/view/settings"
	"net/http"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	return render(r, w, settings.Index(user))
}

func HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	var errors settings.ProfileErrors
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}
	ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(40)),
	}).Validate(&errors)
	if !ok {
		return render(r, w, settings.ProfileForm(params, errors))
	}
	user := getAuthenticatedUser(r)
	user.Account.UserName = params.Username

	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return render(r, w, settings.ProfileForm(params, errors))
}
