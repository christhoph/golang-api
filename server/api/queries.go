package api

import (
	"github.com/graphql-go/graphql"
	"github.com/raunofreiberg/kyrene/server"
	"github.com/raunofreiberg/kyrene/server/api/users"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"getUser": &graphql.Field{
			Type: UserType,
			Description: "Query for getting user logged in",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				jwt := params.Context.Value("jwt").(string)
				_, err := server.ValidateJWT(jwt)

				if err != nil {
					return nil, err
				}

				claims, err := server.ParseToken(jwt)
				userID := int(claims["id"].(float64))
				res, err := users.QueryUserById(userID)

				if err != nil {
					return nil, err
				}

				return res, nil
			},
		},
	},
})
