package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/raunofreiberg/kyrene/server"
	"github.com/raunofreiberg/kyrene/server/api"
	"github.com/raunofreiberg/kyrene/server/database"
)

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    api.RootQuery,
	Mutation: api.RootMutation,
})

func main() {
	h := handler.New(&handler.Config{
		Schema:   &Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	development := os.Getenv("ENV") == "development"

	database.Init()

	if development {
		http.Handle("/", http.FileServer(http.Dir("./client")))
	} // only serve static files in development via this server. Nginx is used in production instead

	http.Handle("/graphql", server.PassJwtContext(h))
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}
