package models

type RecipeComment struct {
	ID     string `json:"id"`
	Comment string `json:"comment"`
	User    User   `json:"user"`
}

type Ingredient struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


type Step struct {
	ID          string `json:"id"`
	StepNumber  int    `json:"step_number"`
	Description string `json:"description"`
}


type FullRecipe struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	FeaturedImage string         `json:"featured_image"`
	CreatedAt     string         `json:"created_at"`
	Difficulty      string  	 `json:"difficulty"`
	PrepTimeMinutes int     	 `json:"prep_time_minutes"`
	CookTimeMinutes int     	 `json:"cook_time_minutes"`
	User          User           `json:"user"`
	Category      Category       `json:"category"`
	Ingredients   []Ingredient   `json:"ingredients"`
	Comments      []RecipeComment `json:"recipe_comments"`
	Steps         []Step          `json:"recipe_steps"`
}
