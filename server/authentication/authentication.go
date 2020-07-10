package authentication

import (
	"errors"
	"time"
	"math/rand"
	"github.com/raunofreiberg/kyrene/server/api/users"
	"github.com/raunofreiberg/kyrene/server/model"
)

type Picture struct {
	Height		int						`json:"height"`
	Url				string				`json:"url"`
	Width			int						`json:"width"`
}

type PictureData struct {
	Data			Picture				`json:"data"`
}

type FacebookMe struct{
	Id				string				`json:"id"`
	Name			string				`json:"name"`
	Email			string				`json:"email"`
	Picture		PictureData		`json:"picture"`
}

var facebookSalt = "";

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func RegisterUser(name string, email string, password string, isFacebook bool) (interface{}, error) {
	hashedPassword, err := HashPassword(password)

	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())

	user, err := users.CreateUser(name, email, hashedPassword, isFacebook)
	if err != nil {
		return nil, err
	}

	signedToken, err := GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return model.Token{
		Token: signedToken,
	}, nil
}

func LoginUser(email string, password string) (interface{}, error) {
	queriedUser, err := users.QueryUser(email)
	if err != nil {
		return nil, errors.New("User not found")
	}

	isAuthenticated, err := users.IsAuthenticated(email, []byte(password))
	if err != nil {
		return nil, err
	}

	if isAuthenticated {
		signedToken, err := GenerateJWT(queriedUser)

		if err != nil {
			return nil, err
		}

		return model.Token{
			Token: signedToken,
		}, nil
	}

	return nil, nil
}

func LoginAdminUser(email string, password string) (interface{}, error) {
	queriedUser, err := users.QueryUser(email)
	if err != nil {
		return nil, errors.New("User not found")
	}

	isAuthenticated, err := users.IsAuthenticated(email, []byte(password))
	if err != nil {
		return nil, err
	}

	if isAuthenticated {
		_, err := users.QueryRole(queriedUser.(model.User).Id)
		if err != nil {
			return nil, errors.New("User has no existing role")
		}

		signedToken, err := GenerateJWT(queriedUser)

		if err != nil {
			return nil, err
		}

		return model.Token{
			Token: signedToken,
		}, nil
	}

	return nil, nil
}

func FacebookAuth(accessToken string, userId string) (interface{}, error) {
	facebookMe := new(FacebookMe)
    GetHttpJson("https://graph.facebook.com/"+ userId + "?fields=picture,birthday,id,name,email&access_token=" + accessToken, facebookMe)

	if facebookMe.Id == "" || facebookMe.Email == "" || facebookMe.Name == "" {
		return nil, errors.New("You are not authorized")
	}

	queriedUser, err := users.QueryUser(facebookMe.Email)

	if err != nil {
		token, err := RegisterUser(facebookMe.Name, facebookMe.Email, facebookMe.Id + facebookSalt, true)

		if err != nil {
			return nil, err
		}

		return token, nil
	}

	if queriedUser.(model.User).IsFacebook {
		token, err := LoginUser(facebookMe.Email, facebookMe.Id + facebookSalt)

		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, errors.New("Ya cuentas con una cuenta de email regular, intenta con esa")
}
