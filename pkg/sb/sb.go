package sb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/nedpals/supabase-go"
)

var (
	Client *supabase.Client
	apiKey string
)

func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	if sbHost == "" {
		return errors.New("supabase url is required")
	}
	sbSecret := os.Getenv("SUPABASE_SECRET")
	if sbSecret == "" {
		return errors.New("supabase secret is required")
	}
	Client = supabase.CreateClient(sbHost, sbSecret)
	apiKey = sbSecret
	return nil
}

const authEndpoint = "auth/v1"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (err *ErrorResponse) Error() string {
	return err.Message
}

// this func simulate the ResetPasswordForEmail from the supabase package but with a redirectTo
func MyResetPasswordForEmail(ctx context.Context, email string, redirectTo string) error {
	reqBody, _ := json.Marshal(map[string]string{"email": email})
	reqURL := fmt.Sprintf("%s/%s/recover", Client.BaseURL, authEndpoint)
	if len(redirectTo) > 0 {
		reqURL += fmt.Sprintf("?redirect_to=%s", redirectTo)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if err = sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}

func sendRequest(req *http.Request, v interface{}) error {
	var errRes ErrorResponse
	hasCustomError, err := sendCustomRequest(req, v, &errRes)

	if err != nil {
		return err
	} else if hasCustomError {
		return &errRes
	}

	return nil
}

func sendCustomRequest(req *http.Request, successValue interface{}, errorValue interface{}) (bool, error) {
	req.Header.Set("apikey", apiKey)
	res, err := Client.HTTPClient.Do(req)
	if err != nil {
		return true, err
	}

	defer res.Body.Close()
	statusOK := res.StatusCode >= http.StatusOK && res.StatusCode < 300
	if !statusOK {
		if err = json.NewDecoder(res.Body).Decode(&errorValue); err == nil {
			return true, nil
		}

		return false, fmt.Errorf("unknown, status code: %d", res.StatusCode)
	} else if res.StatusCode != http.StatusNoContent {
		if err = json.NewDecoder(res.Body).Decode(&successValue); err != nil {
			return false, err
		}
	}

	return false, nil
}
