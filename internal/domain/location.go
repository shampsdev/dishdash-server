package domain

import (
	geo "github.com/kellydunn/golang-geo"
)

func ParsePoint(s string, point *geo.Point) error {
	return point.UnmarshalBinary([]byte(s))
}

func Point2String(p *geo.Point) string {
	bytes, _ := p.MarshalJSON()
	return string(bytes)
}
