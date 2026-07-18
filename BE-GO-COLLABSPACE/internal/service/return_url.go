package service

import (
	"net/url"
	"strings"

	"Go-CollabSpace/internal/common/apperror"
)

func buildAllowedReturnURL(rawReturnURL string, allowedReturnURLs []string, token string) (string, error) {
	returnURL, err := validateAllowedReturnURL(rawReturnURL, allowedReturnURLs)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(returnURL)
	if err != nil {
		return "", apperror.ErrReturnURLNotAllowed
	}
	query := parsed.Query()
	query.Set("token", token)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func validateAllowedReturnURL(rawReturnURL string, allowedReturnURLs []string) (string, error) {
	returnURL, err := normalizeReturnURL(rawReturnURL)
	if err != nil {
		return "", apperror.ErrReturnURLNotAllowed
	}

	for _, allowedURL := range allowedReturnURLs {
		normalizedAllowedURL, err := normalizeReturnURL(allowedURL)
		if err != nil {
			continue
		}
		if returnURL == normalizedAllowedURL {
			return returnURL, nil
		}
	}

	return "", apperror.ErrReturnURLNotAllowed
}

func normalizeReturnURL(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", apperror.ErrReturnURLNotAllowed
	}
	if parsed.Host == "" || parsed.User != nil || parsed.Fragment != "" {
		return "", apperror.ErrReturnURLNotAllowed
	}

	parsed.Scheme = strings.ToLower(parsed.Scheme)
	parsed.Host = strings.ToLower(parsed.Host)
	return parsed.String(), nil
}
