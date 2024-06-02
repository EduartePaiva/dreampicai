package handler

import (
	"dreampicai/db"
	"dreampicai/types"
	"dreampicai/view/generate"
	"fmt"
	"log/slog"
	"net/http"

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
		Images: images,
	}
	fmt.Println(images[0].Status)
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	prompt := "red sport car"
	user := getAuthenticatedUser(r)
	img := types.Image{
		Prompt: prompt,
		UserID: user.ID,
		Status: types.ImageStatusPending,
	}

	if err := db.CreateImage(&img); err != nil {
		return err
	}

	return render(r, w, generate.GalleryImage(img))
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
