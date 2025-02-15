package domain

type Collection struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Places      []*Place `json:"places"`
}
