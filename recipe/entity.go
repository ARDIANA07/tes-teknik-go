package recipe

type Recipe struct {
	ID          string   `json:"id"`
	Name        string   `json:"title" type:"varchar(50)"`
	Description string   `json:"description" type:"text"`
	Ingredients []string `json:"ingredients"`
	Category    string   `json:"category"`
}
