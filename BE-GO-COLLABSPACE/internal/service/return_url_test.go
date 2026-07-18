package service

import (
	"errors"
	"net/url"
	"testing"

	"Go-CollabSpace/internal/common/apperror"
)

func TestBuildAllowedReturnURL(t *testing.T) {
	allowed := []string{
		"https://app.example.com/reset-password",
		"http://localhost:3000/reset-password",
	}

	resetURL, err := buildAllowedReturnURL("HTTPS://APP.EXAMPLE.COM/reset-password", allowed, "abc 123")
	if err != nil {
		t.Fatalf("expected allowed return URL, got error: %v", err)
	}

	parsed, err := url.Parse(resetURL)
	if err != nil {
		t.Fatalf("expected valid URL, got error: %v", err)
	}
	if parsed.Scheme != "https" || parsed.Host != "app.example.com" || parsed.Path != "/reset-password" {
		t.Fatalf("unexpected reset URL: %s", resetURL)
	}
	if parsed.Query().Get("token") != "abc 123" {
		t.Fatalf("expected token query param, got %q", parsed.Query().Get("token"))
	}
}

func TestBuildAllowedReturnURLPreservesAllowedQuery(t *testing.T) {
	allowed := []string{"https://app.example.com/reset-password?source=email"}

	resetURL, err := buildAllowedReturnURL("https://app.example.com/reset-password?source=email", allowed, "token")
	if err != nil {
		t.Fatalf("expected allowed return URL, got error: %v", err)
	}

	parsed, err := url.Parse(resetURL)
	if err != nil {
		t.Fatalf("expected valid URL, got error: %v", err)
	}
	if parsed.Query().Get("source") != "email" {
		t.Fatalf("expected existing query param to be preserved, got %q", parsed.Query().Get("source"))
	}
	if parsed.Query().Get("token") != "token" {
		t.Fatalf("expected token query param, got %q", parsed.Query().Get("token"))
	}
}

func TestBuildAllowedReturnURLRejectsUnsafeURLs(t *testing.T) {
	allowed := []string{"https://app.example.com/reset-password"}
	tests := []string{
		"https://app.example.com.evil/reset-password",
		"https://evil.example.com/reset-password",
		"https://app.example.com/other",
		"javascript:alert(1)",
		"https://app.example.com/reset-password#token=steal",
		"https://app.example.com@evil.example.com/reset-password",
	}

	for _, testURL := range tests {
		t.Run(testURL, func(t *testing.T) {
			_, err := buildAllowedReturnURL(testURL, allowed, "token")
			if !errors.Is(err, apperror.ErrReturnURLNotAllowed) {
				t.Fatalf("expected ErrReturnURLNotAllowed, got %v", err)
			}
		})
	}
}
