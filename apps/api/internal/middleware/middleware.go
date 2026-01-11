package middleware

import (
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(next)
}
