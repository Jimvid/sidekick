package user

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
)

func GetUserId(r *http.Request) (string, error) {
	claims, ok := clerk.SessionClaimsFromContext(r.Context())
	if !ok {
		slog.Warn("No session claims in context")
		return "", fmt.Errorf("No session claims in context")
	}

	return claims.Subject, nil
}
