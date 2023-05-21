package ingredient

type Ingredient struct {
	ID          string `json:"id"`
	Name        string `json:"name" type:"varchar(255)"`
	JumlahBahan string `json:"jumlah_bahan" type:"varchar(20)"`
}
