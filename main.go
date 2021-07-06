package main

import (
	"SpaceshipAsteroids/server/communicator"
	"SpaceshipAsteroids/server/handler"
	"log"
	"net/http"

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
	r.HandleFunc("/topsecret_split/{satteliteName}", nil).Methods(http.MethodPost)
	r.HandleFunc("/topsecret_split/{satteliteName}", nil).Methods(http.MethodGet)

	log.Println("Listening on port 8080")
	return http.ListenAndServe(":8080", r)
}
