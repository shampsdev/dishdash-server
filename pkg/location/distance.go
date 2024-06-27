package location

import (
	"math"

	"dishdash.ru/internal/domain"
)

func GetDistance(from, to domain.Coordinate) int64 {
	const R = 6371
	dLat := deg2rad(to.Lat - from.Lat)
	dLon := deg2rad(to.Lon - from.Lon)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(from.Lat))*math.Cos(deg2rad(to.Lat))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return int64(d * 1000)
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
