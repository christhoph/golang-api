package server

import (
	"context"
	"errors"
	"net/http"
	"regexp"
)

func PassJwtContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		bearerRegex, _ := regexp.Compile("(?:Bearer *)([^ ]+)(?: *)")
		bearerRegexMatches := bearerRegex.FindStringSubmatch(authorizationHeader)
		// We still want requests to be passed onto GraphQL queries/mutations
		// for authentication so we set the JWT to an empty string if it doesn't exist.
		ctx := context.WithValue(r.Context(), "jwt", "")

		if len(bearerRegexMatches) != 0 {
			jwtToken := bearerRegexMatches[1]
			ctx = context.WithValue(r.Context(), "jwt", jwtToken)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper middleware func for wrapping certain queries in auth verification
// E.g: RequireAuth(jwt, someQuery)
func RequireAuth(jwt string, callback func() (interface{}, error)) (interface{}, error) {
	isAuthorized, err := ValidateJWT(jwt)

	if err != nil {
		return nil, err
	}

	if isAuthorized {
		res, err := callback()

		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, errors.New("Not authorized")
}
