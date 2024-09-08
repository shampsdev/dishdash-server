package domain

import (
	"math"
)

const EARTH_RADIUS = 6371.0

type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func (c *Coordinate) GreatCircleDistance(c2 *Coordinate) float64 {
	dLat := (c2.Lat - c.Lat) * (math.Pi / 180.0)
	dLon := (c2.Lon - c.Lon) * (math.Pi / 180.0)

	lat1 := c.Lat * (math.Pi / 180.0)
	lat2 := c2.Lat * (math.Pi / 180.0)

	a1 := math.Sin(dLat/2) * math.Sin(dLat/2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2
	circle := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * circle
}
