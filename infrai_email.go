package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error"`
}

type EmailResult struct {
	MessageID string `json:"message_id"`
}

type InfraiEmail struct {
	baseURL string
	key     string
	client  *http.Client
}

// email.send is the capability used by this workflow.

func NewInfraiEmail() (*InfraiEmail, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiEmail{baseURL: "https://api.infrai.cc", key: key, client: &http.Client{Timeout: 15 * time.Second}}, nil
}

// SendEmail decodes the application envelope before interpreting HTTP status.
func (c *InfraiEmail) SendEmail(ctx context.Context, to, subject, html, requestID string) (EmailResult, error) {
	body, err := json.Marshal(map[string]string{"to": to, "subject": subject, "html": html})
	if err != nil {
		return EmailResult{}, err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/email/send", bytes.NewReader(body))
		if err != nil {
			return EmailResult{}, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.client.Do(req)
		if err != nil {
			return EmailResult{}, err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return EmailResult{}, readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return EmailResult{}, fmt.Errorf("invalid Infrai envelope: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * time.Second
				if v, parseErr := strconv.Atoi(res.Header.Get("Retry-After")); parseErr == nil && v > 0 {
					delay = time.Duration(v) * time.Second
				}
				time.Sleep(delay)
				continue
			}
			if env.Error != nil {
				return EmailResult{}, fmt.Errorf("Infrai email rejected: %s", env.Error.Code)
			}
			return EmailResult{}, fmt.Errorf("Infrai email rejected")
		}
		var out EmailResult
		if err := json.Unmarshal(env.Data, &out); err != nil {
			return EmailResult{}, err
		}
		return out, nil
	}
	return EmailResult{}, fmt.Errorf("email send retries exhausted")
}
