package main

import (
	"context"
	"dreampicai/db"
	"dreampicai/handler"
	"dreampicai/pkg/sb"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

//go:embed public
var FS embed.FS

func foo() {

	ctx := context.Background()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		// handle error
		slog.Error("api key error")
		return
	}

	// https://replicate.com/stability-ai/stable-diffusion
	version := "ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4"

	input := replicate.PredictionInput{
		"prompt": "an astronaut riding a horse on mars, hd, dramatic lighting",
	}

	webhook := replicate.Webhook{
		URL:    "https://example.com/webhook",
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, err := r8.Run(ctx, version, input, &webhook)
	if err != nil {
		// handle error
	}
	fmt.Println("output: ", output)

	// Create a prediction
	prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		// handle error
	}

	// Wait for the prediction to finish
	err = r8.Wait(ctx, prediction)
	if err != nil {
		// handle error
	}
	fmt.Println("output: ", output)

	// The `Run` method is a convenience method that
	// creates a prediction, waits for it to finish, and returns the output.
	// If you want a reference to the prediction, you can call `CreatePrediction`,
	// call `Wait` on the prediction, and access its `Output` field.
	prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		// handle error
	}

	// Wait for the prediction to finish
	err = r8.Wait(ctx, prediction)
	if err != nil {
		// handle error
	}
	fmt.Println("output: ", prediction.Output)
}

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	foo()
	return

	router := chi.NewMux()
	router.Use(handler.WithUser)

	router.Handle("/*", public())
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth, handler.WithoutAccountSetup)
		auth.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		auth.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))
	})

	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAuth, handler.WithAccountSetup)
		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.Make(handler.HandleSettingsUsernameUpdate))
		auth.Post("/settings/account/password", handler.Make(handler.HandleResetPasswordCreate))
		auth.Get("/auth/reset-password", handler.Make(handler.HandleResetPasswordIndex))
		auth.Put("/auth/reset-password", handler.Make(handler.HandleResetPasswordUpdate))

		auth.Get("/generate", handler.Make(handler.HandleGenerateIndex))
		auth.Post("/generate", handler.Make(handler.HandleGenerateCreate))
		auth.Get("/generate/image/status/{id}", handler.Make(handler.HandleGenerateImageStatus))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
