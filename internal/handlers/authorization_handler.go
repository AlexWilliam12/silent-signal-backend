package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/auth"
)

func HandleAuthorization(r *http.Request) (*auth.CustomClaims, error) {

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return nil, fmt.Errorf("unauthorized request")
	}

	if !strings.Contains(authorization, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization request")
	}

	token := strings.Replace(authorization, "Bearer ", "", 1)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("unauthorized request")
	}

	return claims, nil
}
