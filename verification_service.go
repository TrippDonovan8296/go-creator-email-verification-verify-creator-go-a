package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Signup struct {
	Email string
	Name  string
}
type Verification struct {
	Email     string
	Link      string
	MessageID string
}

func newToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func BuildVerificationLink(baseURL, email, token string) string {
	return fmt.Sprintf("%s/verify?email=%s&token=%s", baseURL, email, token)
}

func SendVerification(ctx context.Context, mailer *InfraiEmail, signup Signup, baseURL string) (Verification, error) {
	if signup.Email == "" {
		return Verification{}, fmt.Errorf("email is required")
	}
	token, err := newToken()
	if err != nil {
		return Verification{}, err
	}
	link := BuildVerificationLink(baseURL, signup.Email, token)
	html := fmt.Sprintf("<p>Hello %s,</p><p>Confirm your creator account:</p><p><a href=\"%s\">Verify email</a></p>", signup.Name, link)
	result, err := mailer.SendEmail(ctx, signup.Email, "Verify your creator account", html, "signup-"+token)
	if err != nil {
		return Verification{}, err
	}
	return Verification{Email: signup.Email, Link: link, MessageID: result.MessageID}, nil
}
