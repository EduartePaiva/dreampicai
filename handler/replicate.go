package handler

import (
	"context"
	"database/sql"
	"dreampicai/db"
	"dreampicai/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// "input": {
//     "apply_watermark": true,
//     "guidance_scale": 2,
//     "height": 1024,
//     "negative_prompt": "CGI, Unreal, Airbrushed, Digital",
//     "num_inference_steps": 5,
//     "num_outputs": 1,
//     "prompt": "a nice body woman in a red dress",
//     "scheduler": "DPM++SDE",
//     "width": 1024
//   }

type ReplicateResp struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
	Input  struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
}

const (
	processing = "processing"
	succeeded  = "succeeded"
)

func HandleReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResp

	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Status == processing {
		return nil
	}
	if resp.Status != succeeded {
		return fmt.Errorf("replicate callback response with a non ok status: %s", resp.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		return fmt.Errorf("replicate callback invalid batchID: %s", err)
	}

	images, err := db.GetImagesByBatchID(batchID)
	if err != nil {
		return fmt.Errorf("replicate callback failed to find image with batchID %s: %s", batchID, err)
	}

	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback un-equal images compared to replicate output")
	}
	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		for i, imageURL := range resp.Output {
			images[i].Status = types.ImageStatusCompleted
			images[i].ImageLocation = imageURL
			images[i].Prompt = resp.Input.Prompt
			if err := db.UpdateImage(&images[i]); err != nil {
				return err
			}
		}
		return nil
	})
	// fmt.Printf("%+v\n", resp)
	return err
}
