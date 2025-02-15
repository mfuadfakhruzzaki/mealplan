package spoonacular

type MealPlanResponse struct {
	Meals     []Meal `json:"meals"`
	Nutrients struct {
		Calories      float64 `json:"calories"`
		Protein       float64 `json:"protein"`
		Fat           float64 `json:"fat"`
		Carbohydrates float64 `json:"carbohydrates"`
	} `json:"nutrients"`
}

type Meal struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	ImageType      string `json:"imageType"`
	ReadyInMinutes int    `json:"readyInMinutes"`
	Servings       int    `json:"servings"`
	SourceURL      string `json:"sourceUrl"`
}

type RecipeInformation struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Servings  int    `json:"servings"`
	Nutrition struct {
		Nutrients []Nutrient `json:"nutrients"`
	} `json:"nutrition"`
}

type Nutrient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Unit   string  `json:"unit"`
}
