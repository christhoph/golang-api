package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-test/deep"
	"github.com/raunofreiberg/kyrene/server"
)

type ResponseMock struct {
	Name string
}

func TestJwtMiddleware_WithToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/graphql", nil)
	req.Header.Add("Authorization", "Bearer ey28718921th")

	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, ok := r.Context().Value("jwt").(string)

		if !ok {
			t.Errorf("jwt not in request context: got %q", jwt)
		}

		if len(jwt) == 0 {
			t.Errorf("jwt should not be empty when it's passed")
		}
	})

	rr := httptest.NewRecorder()
	handler := server.PassJwtContext(testHandler)
	handler.ServeHTTP(rr, req)
}

// Don't provide a Authorization header, JWT should be a empty string since we
// pass the context to the GraphQL resolver func that handles authentication.
func TestJwtMiddleware_WithoutToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/graphql", nil)

	if err != nil {
		t.Fatal(err)
	}

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, ok := r.Context().Value("jwt").(string)

		if !ok {
			t.Errorf("jwt not in request context: got %q", jwt)
		}

		if len(jwt) != 0 {
			t.Errorf("jwt should be empty when it's not passed")
		}
	})

	rr := httptest.NewRecorder()
	handler := server.PassJwtContext(testHandler)
	handler.ServeHTTP(rr, req)
}

func TestRequireAuthMiddleware_ValidToken(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, _ := token.SignedString(server.JwtSecret)
	resMock := ResponseMock{"Doge"}
	callbackMock := func() (interface{}, error) {
		return resMock, nil
	}

	res, err := server.RequireAuth(tokenString, callbackMock)

	if err != nil {
		t.Error(err)
	}

	if diff := deep.Equal(res, resMock); diff != nil {
		t.Error(diff)
	}
}

func TestRequireAuthMiddleware_InvalidToken(t *testing.T) {
	invalidToken := "123123"
	callbackMock := func() (interface{}, error) {
		return ResponseMock{}, nil
	}

	_, err := server.RequireAuth(invalidToken, callbackMock)

	if diff := deep.Equal(err.Error(), "Invalid token"); diff != nil {
		t.Error(diff)
	}
}
