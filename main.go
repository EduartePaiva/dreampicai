package main

import (
	"dreampicai/db"
	"dreampicai/handler"
	"dreampicai/pkg/sb"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

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
	router.Post("/replicate/callback/{userID}/{batchID}", handler.Make(handler.HandleReplicateCallback))

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
