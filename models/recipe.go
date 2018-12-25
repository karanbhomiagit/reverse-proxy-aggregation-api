package models

type Imagelink string

type Ingredient struct {
	Name      string    `json:"name"`
	ImageLink Imagelink `json:"imageLink"`
}

type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Headline    string       `json:"headline"`
	Description string       `json:"description"`
	Difficulty  int          `json:"difficulty"`
	PrepTime    string       `json:"prepTime"`
	ImageLink   Imagelink    `json:"imageLink"`
	Ingredients []Ingredient `json:"ingredients"`
}
