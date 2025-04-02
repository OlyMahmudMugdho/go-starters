package server

import (
	"clean-rest-api/handlers"
	cleanapi "clean-rest-api/pkg/clean-api"
	"clean-rest-api/types"
	"errors"
	"log"
	"net/http"
)

func Run() {
	c := cleanapi.New()

	c.Get("/", func(ctx *types.Context) error {
		msgHandler := handlers.NewMessageHandler()
		ctx.Json(msgHandler.GetAllMessages())
		return errors.New("demo error")
	})

	log.Println("server is running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Router))
}
