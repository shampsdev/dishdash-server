package domain

import (
	"encoding/json"

	geo "github.com/kellydunn/golang-geo"
)

func ParsePoint(s string) (*geo.Point, error) {
	p := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	err := json.Unmarshal([]byte(s), &p)

	return geo.NewPoint(p.Latitude, p.Longitude), err
}

func Point2String(point *geo.Point) string {
	p := struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{
		Latitude:  point.Lat(),
		Longitude: point.Lng(),
	}

	bytes, _ := json.Marshal(p)
	return string(bytes)
}
