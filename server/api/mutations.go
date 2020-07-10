package api

import (
	"github.com/graphql-go/graphql"
	"github.com/raunofreiberg/kyrene/server/authentication"
)

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"registerUser": &graphql.Field{
			Type: LoginType,
			Description: "Register user",
			Args:	graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				name := params.Args["name"].(string)
				email := params.Args["email"].(string)
				password := params.Args["password"].(string)
				token, err := authentication.RegisterUser(name, email, password, false)

				if err != nil {
					return nil, err
				}

				return token, nil
			},
		},
		"loginUser": &graphql.Field{
			Type: LoginType,
			Description: "Login user",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				email := params.Args["email"].(string)
				password := params.Args["password"].(string)
				token, err := authentication.LoginUser(email, password)

				if err != nil {
					return nil, err
				}

				return token, nil
			},
		},
	},
})