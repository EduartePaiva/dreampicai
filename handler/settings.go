package handler

import (
	"dreampicai/view/settings"
	"net/http"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	return render(r, w, settings.Index(user))
}
