package handler

import (
	"dreampicai/db"
	"dreampicai/pkg/kit/validate"
	"dreampicai/view/generate"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images:     images,
		FormParams: generate.FormParams{Amount: 1},
	}
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	// user := getAuthenticatedUser(r)

	amount, _ := strconv.Atoi(r.FormValue("amount"))
	params := generate.FormParams{
		Prompt: r.FormValue("prompt"),
		Amount: amount,
	}

	var errors generate.FormErrors
	ok := validate.New(params, validate.Fields{
		"Prompt": validate.Rules(validate.Min(10), validate.Max(100)),
		"Amount": validate.Rules(validate.OnlyTheseNumbers([]int{1, 2, 4, 8})),
	}).Validate(&errors)
	if !ok {
		return render(r, w, generate.Form(params, errors))
	}

	return nil
	// img := types.Image{
	// 	Prompt: prompt,
	// 	UserID: user.ID,
	// 	Status: types.ImageStatusPending,
	// }

	// if err := db.CreateImage(&img); err != nil {
	// 	return err
	// }

	// return render(r, w, generate.GalleryImage(img))
}
func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	//fetch from db
	image, err := db.GetImageByID(id)
	if err != nil {
		return err
	}
	slog.Info("checking image status", "id", id)
	slog.Info("checking image status", "status", image.Status)
	return render(r, w, generate.GalleryImage(image))
}
