package handler

import (
	"dreampicai/types"
	"dreampicai/view/generate"
	"net/http"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	data := generate.ViewData{
		Images: []types.Image{},
	}
	return render(r, w, generate.Index(data))
}
