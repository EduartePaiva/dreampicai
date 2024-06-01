package handler

import (
	"dreampicai/view/generate"
	"net/http"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, generate.Index())
}
