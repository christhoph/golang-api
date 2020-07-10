package model

import (
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id					int					`json:"id,omitempty" sql:",pk"`
	Name				string			`json:"name,omitempty"`
	Email				string			`json:"email,omitempty"`
	Password		[]uint8			`json:"-"`
	Role				*Role				`json:"role,omitempty"`
	IsFacebook	bool				`json:"isFacebook"`
	IsGoogle		bool				`json:"isGoogle"`
	jwt.StandardClaims
}
