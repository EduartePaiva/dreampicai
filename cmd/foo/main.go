package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	foo()
}

func foo() {

	ctx := context.Background()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	// https://replicate.com/stability-ai/stable-diffusion
	version := "bea09cf018e513cef0841719559ea86d2299e05448633ac8fe270b5d5cd6777e"

	input := replicate.PredictionInput{
		"prompt":              "a nice body woman in a red dress",
		"width":               1024,
		"height":              1024,
		"scheduler":           "DPM++SDE",
		"num_outputs":         1,
		"guidance_scale":      2,
		"apply_watermark":     true,
		"negative_prompt":     "CGI, Unreal, Airbrushed, Digital",
		"num_inference_steps": 5,
	}

	webhook := replicate.Webhook{
		URL:    os.Getenv("REPLICATE_CALLBACK_URL"),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	output, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ", output)

	// Create a prediction
	// prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	// if err != nil {
	// 	// handle error
	// }

	// Wait for the prediction to finish
	// err = r8.Wait(ctx, prediction)
	// if err != nil {
	// 	// handle error
	// }
	// fmt.Println("output: ", output)

	// The `Run` method is a convenience method that
	// creates a prediction, waits for it to finish, and returns the output.
	// If you want a reference to the prediction, you can call `CreatePrediction`,
	// call `Wait` on the prediction, and access its `Output` field.
	// prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	// if err != nil {
	// 	// handle error
	// }

	// // Wait for the prediction to finish
	// err = r8.Wait(ctx, prediction)
	// if err != nil {
	// 	// handle error
	// }
	// fmt.Println("output: ", prediction.Output)
}
