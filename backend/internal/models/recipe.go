package models



type Recipe struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	FeaturedImage   string  `json:"featured_image"`
	Difficulty      string  `json:"difficulty"`
	PrepTimeMinutes int     `json:"prep_time_minutes"`
	CookTimeMinutes int     `json:"cook_time_minutes"`
	UserID          string  `json:"user_id"`
	CategoryID      string  `json:"category_id"`
	IsPaid          bool    `json:"is_paid"`
	Price           float64 `json:"price"`
	CreatedAt       string  `json:"created_at"`   
	UpdatedAt       string  `json:"updated_at"`   
}
