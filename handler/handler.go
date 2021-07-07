package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"SpaceshipAsteroids/server/communicator"
)

type Communicator interface {
	GetSpaceshipCoordinates(satellites ...communicator.Satellite) (communicator.GetSpaceshipCoordinatesResponse, error)
	GetSpaceshipCoordinatesByOne(satellite communicator.Satellite) (communicator.GetSpaceshipCoordinatesResponse, error)
}

type Handler struct {
	communicator Communicator
}

func NewHandler(communicator Communicator) *Handler {
	return &Handler{communicator: communicator}
}

func (h *Handler) GetSpaceshipCoordinates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Satellites []struct {
				Name     string   `json:"name"`
				Distance float32  `json:"distance"`
				Message  []string `json:"message"`
			} `json:"satellites"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			prepareErr(w, http.StatusUnprocessableEntity, fmt.Sprintf("could not process body: %v", err))
			return
		}

		var ss []communicator.Satellite
		for _, s := range body.Satellites {
			ss = append(ss, communicator.Satellite{
				Name:     s.Name,
				Distance: s.Distance,
				Message:  s.Message,
			})
		}

		coordinates, err := h.communicator.GetSpaceshipCoordinates(ss...)
		if err != nil {
			prepareErr(w, http.StatusConflict, err.Error())
			return
		}

		b, err := json.Marshal(coordinates)
		if err != nil {
			prepareErr(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, _ = w.Write(b)
	}
}

func (h *Handler) GetSpaceshipCoordinatesByOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Satellites struct {
				Name     string   `json:"name"`
				Distance float32  `json:"distance"`
				Message  []string `json:"message"`
			} `json:"satellites"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			prepareErr(w, http.StatusUnprocessableEntity, fmt.Sprintf("could not process body: %v", err))
			return
		}

		var ss = communicator.Satellite{
			Name:     body.Satellites.Name,
			Distance: body.Satellites.Distance,
			Message:  body.Satellites.Message,
		}

		coordinates, err := h.communicator.GetSpaceshipCoordinatesByOne(ss)
		if err != nil {
			prepareErr(w, http.StatusConflict, err.Error())
			return
		}

		b, err := json.Marshal(coordinates)
		if err != nil {
			prepareErr(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, _ = w.Write(b)
	}
}

func prepareErr(w http.ResponseWriter, code int, message string) {
	responseBody := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Message: message,
		Code:    code,
	}

	b, err := json.Marshal(responseBody)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(code)
	_, _ = w.Write(b)
}
