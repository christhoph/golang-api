package api

import (
	"github.com/graphql-go/graphql"
	"github.com/raunofreiberg/kyrene/server/model"
	"github.com/raunofreiberg/kyrene/server/api/users"
)

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"role": &graphql.Field{
			Type: RoleType,
			Description: "query user role",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := params.Source.(model.User).Id
				res, err := users.QueryRole(userId)

				if err != nil {
					return nil, err
				}

				return res, nil
			},
		},
	},
})

var RoleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RoleType",
	Fields: graphql.Fields {
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"userId": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"isAdmin": &graphql.Field{
			Type: graphql.Boolean,
		},
		"isGod": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var LoginType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Login",
	Fields: graphql.Fields {
		"token": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var SuccessType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Success",
	Fields: graphql.Fields{
		"ok": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})
