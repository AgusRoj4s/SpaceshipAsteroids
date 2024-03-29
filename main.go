package main

import (
	"SpaceshipAsteroids/server/communicator"
	"SpaceshipAsteroids/server/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := communicator.NewCommunicator()

	h := handler.NewHandler(c)

	r := mux.NewRouter()
	r.HandleFunc("/topsecret", h.GetSpaceshipCoordinates()).Methods(http.MethodPost)
	r.HandleFunc("/topsecret_split/{satteliteName}", h.GetSpaceshipCoordinatesByOne()).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Listening on port " + port)
	return http.ListenAndServe(":"+port, r)
}
