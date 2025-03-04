package geo

import (
	"sort"

	"dishdash.ru/pkg/domain"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

func postgisDistance(coord1, coord2 domain.Coordinate) float64 {
	p1 := orb.Point{coord1.Lon, coord1.Lat}
	p2 := orb.Point{coord2.Lon, coord2.Lat}
	return geo.DistanceHaversine(p1, p2)
}

func SortPlacesByDistance(places []*domain.Place, refPoint domain.Coordinate) {
	sort.Slice(places, func(i, j int) bool {
		distI := postgisDistance(places[i].Location, refPoint)
		distJ := postgisDistance(places[j].Location, refPoint)
		return distI < distJ
	})
}
