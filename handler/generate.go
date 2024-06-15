package handler

import (
	"context"
	"database/sql"
	"dreampicai/db"
	"dreampicai/pkg/kit/validate"
	"dreampicai/types"
	"dreampicai/view/generate"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"github.com/uptrace/bun"
)

const creditsPerImage = 2

func HandleGenerateIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	images, err := db.GetImagesByUserID(user.ID)
	if err != nil {
		return err
	}
	data := generate.ViewData{
		Images: images,
		// FormParams: generate.FormParams{Amount: 1},
	}
	return render(r, w, generate.Index(data))
}

func HandleGenerateCreate(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

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

	creditsNeeded := params.Amount * creditsPerImage
	if user.Credits < creditsNeeded {
		errors.Credits = true
		errors.CreditsNeeded = creditsNeeded
		errors.UserCredits = user.Account.Credits
		return render(r, w, generate.Form(params, errors))
	}
	return nil

	batchID := uuid.New()
	genParams := GenerateImagesParams{
		Prompt:  params.Prompt,
		Amount:  params.Amount,
		UserID:  user.UserID,
		BatchID: batchID,
	}
	if err := generateImage(r.Context(), genParams); err != nil {
		return err
	}

	err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for range params.Amount {
			img := types.Image{
				UserID:  user.ID,
				Status:  types.ImageStatusPending,
				BatchID: batchID,
			}
			if err := db.CreateImage(&img); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return hxRedirect(w, r, "/generate")
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

type GenerateImagesParams struct {
	Prompt  string
	Amount  int
	BatchID uuid.UUID
	UserID  uuid.UUID
}

func generateImage(ctx context.Context, params GenerateImagesParams) error {
	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		return err
	}

	// https://replicate.com/stability-ai/stable-diffusion
	version := "bea09cf018e513cef0841719559ea86d2299e05448633ac8fe270b5d5cd6777e"

	// this must be a aspect ratio of 2/3
	input := replicate.PredictionInput{
		"prompt":              params.Prompt,
		"width":               640,
		"height":              960,
		"scheduler":           "DPM++SDE",
		"num_outputs":         params.Amount,
		"guidance_scale":      2,
		"apply_watermark":     true,
		"negative_prompt":     "CGI, Unreal, Airbrushed, Digital",
		"num_inference_steps": 5,
	}

	webhook := replicate.Webhook{
		URL:    fmt.Sprintf("%s/%s/%s", os.Getenv("REPLICATE_CALLBACK_URL"), params.UserID, params.BatchID),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	_, err = r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		return err
	}
	return nil
}
