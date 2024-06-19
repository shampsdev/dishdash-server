package domain

type Card struct {
	ID               int64
	Title            string
	ShortDescription string
	Description      string
	Image            string
	Location         Coordinate
	Address          string
	Price            int
	Tags             []*Tag
}
