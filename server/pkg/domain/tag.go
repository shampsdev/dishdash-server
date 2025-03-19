package domain

type Tag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Visible bool   `json:"visible"`
	Order   int    `json:"order"`
	// Places with excluded tag wont be shown in swipes
	Excluded bool `json:"excluded"`
}
