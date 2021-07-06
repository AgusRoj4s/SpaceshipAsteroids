package communicator

import (
	"errors"
	"sort"
	"strings"
)

type Communicator struct {
	satellites satellites
}

func NewCommunicator() *Communicator {
	defaultSatellitesCoordinates := []struct {
		name        string
		coordinates []float32
	}{
		{
			name:        "Kenobi",
			coordinates: []float32{-500, -200},
		},
		{
			name:        "Skywalker",
			coordinates: []float32{100, -100},
		},
		{
			name:        "Sato",
			coordinates: []float32{500, 100},
		},
	}

	ss := make(satellites)
	for _, s := range defaultSatellitesCoordinates {
		ss[s.name] = newSatellite(s.name, s.coordinates[0], s.coordinates[1])
	}

	return &Communicator{
		satellites: ss,
	}
}

type satellites map[string]*satellite

func (ss satellites) getSatellite(name string) (satellite, error) {
	s, exist := ss[name]
	if !exist {
		return satellite{}, errors.New("could not find satellite")
	}

	return *s, nil
}

type satellite struct {
	Name       string
	Coordinate []float32
}

func newSatellite(name string, x, y float32) *satellite {
	return &satellite{
		Name:       name,
		Coordinate: []float32{x, y},
	}
}

func (s *satellite) getX() float32 {
	return s.Coordinate[0]
}

func (s *satellite) getY() float32 {
	return s.Coordinate[1]
}

type Satellite struct {
	Name     string
	Distance float32
	Message  []string
}

type GetSpaceshipCoordinatesResponse struct {
	Position struct {
		X float32 `json:"x"`
		Y float32 `json:"y"`
	} `json:"position"`
	Message string `json:"message"`
}

func (c *Communicator) GetSpaceshipCoordinates(satellites ...Satellite) (GetSpaceshipCoordinatesResponse, error) {
	if len(satellites) == 3 {
		x, y, err := c.getSpaceshipCoordinatesWithThreeDistances(satellites[0].Distance, satellites[1].Distance, satellites[2].Distance)
		if err != nil {
			return GetSpaceshipCoordinatesResponse{}, err
		}

		return GetSpaceshipCoordinatesResponse{
			Position: struct {
				X float32 `json:"x"`
				Y float32 `json:"y"`
			}{X: x, Y: y},
			Message: getSpaceshipMessage(satellites[0].Message, satellites[1].Message, satellites[2].Message),
		}, nil
	}

	return GetSpaceshipCoordinatesResponse{}, errors.New("could not calculate spaceship's coordinates. not enough information")

}

func (c *Communicator) getSpaceshipCoordinatesWithThreeDistances(d1, d2, d3 float32) (x, y float32, err error) {
	k, _ := c.satellites.getSatellite("Kenobi")
	sk, _ := c.satellites.getSatellite("Skywalker")
	s, _ := c.satellites.getSatellite("Sato")

	x1 := k.getX()
	x2 := sk.getX()
	x3 := s.getX()
	y1 := k.getY()
	y2 := sk.getY()
	y3 := s.getY()

	A := x1 - x2
	B := y1 - y2
	D := x1 - x3
	E := y1 - y3

	T := d1*d1 - x1*x1 - y1*y1
	C := (d2*d2 - x2*x2 - y2*y2) - T
	F := (d3*d3 - x3*x3 - y3*y3) - T

	Mx := (C*E - B*F) / 2
	My := (A*F - D*C) / 2
	M := A*E - D*B

	if M == 0 {
		return 0, 0, errors.New("could not calculate the ship's coordinates")
	}

	x = Mx / M
	y = My / M

	return x, y, nil
}

func (c *Communicator) getSpaceshipCoordinatesWithTwoDistances(d1, d2 float32) (x, y float32, err error) {
	k, _ := c.satellites.getSatellite("Kenobi")
	sk, _ := c.satellites.getSatellite("Skywalker")

	d1 = 100
	d2 = 200

	x1 := k.getX()
	x2 := sk.getX()
	y1 := k.getY()
	y2 := sk.getY()

	Mx := (x1*d2 - x2*d1)
	My := (y1*d2 - y2*d1)
	M := x1*y2 - x2*y1

	if M == 0 {
		return 0, 0, errors.New("could not calculate the ship's coordinates")
	}

	x = Mx / M
	y = My / M

	return x, y, nil
}

func getSpaceshipMessage(messages ...[]string) string {
	words := make(map[string]int, len(messages))
	for _, message := range messages {
		for i, word := range message {
			if word == "" {
				continue
			}

			if _, exist := words[word]; !exist || i < words[word] {
				words[word] = i
			}
		}
	}

	type wordIndex struct {
		word  string
		index int
	}

	wordsSorted := make([]wordIndex, 0, len(words))
	for word, minIndexSeen := range words {
		wi := wordIndex{
			word:  word,
			index: minIndexSeen,
		}

		wordsSorted = append(wordsSorted, wi)
	}

	sort.Slice(wordsSorted, func(i, j int) bool {
		return wordsSorted[i].index < wordsSorted[j].index
	})

	var message []string
	for _, wi := range wordsSorted {
		message = append(message, wi.word)
	}

	return strings.Join(message, " ")
}
