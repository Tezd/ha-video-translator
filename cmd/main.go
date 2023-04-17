package main

import (
	"log"

	"ha-video-translator/cmd/server"
	"ha-video-translator/pkg/service"
)

func main() {
	srv := service.New()
	srv.RegisterHttpHandlers()

	if err := server.New(8080).ListenAndServe(); err != nil {
		log.Fatalf("failed to close server: [%v]", err)
	}
}
