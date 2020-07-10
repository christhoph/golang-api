package authentication

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/raunofreiberg/kyrene/server"
	"github.com/raunofreiberg/kyrene/server/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"encoding/json"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

// TODO: are tokens validated against their expiry date?
func GenerateJWT(user interface{}) (string, error) {
	// 4380 hours = 6 months
	expireToken := time.Now().Add(time.Hour * 4380).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.User{
		Id: user.(model.User).Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	})
	signedToken, err := token.SignedString(server.JwtSecret)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func HashPassword(password string) ([]uint8, error) {
	bytePassword := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func GetHttpJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
