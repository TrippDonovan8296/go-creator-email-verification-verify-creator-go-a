package main

import "testing"

func TestBuildVerificationLink(t *testing.T) {
	tests := []struct{ name, base, email, token, want string }{
		{"creator signup", "https://creator.example", "a@b.test", "abc123", "https://creator.example/verify?email=a@b.test&token=abc123"},
		{"local service", "http://localhost:8080", "owner@example.test", "z9", "http://localhost:8080/verify?email=owner@example.test&token=z9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildVerificationLink(tt.base, tt.email, tt.token); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
