package recipe

type Recipe struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	Category    string   `json:"category"`
}
