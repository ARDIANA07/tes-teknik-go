package ingredient

import "database/sql"

type IngredientRepository interface {
	GetAll() ([]Ingredient, error)
	GetByID(id string) (*Ingredient, error)
	Create(ingredient *Ingredient) error
	Update(id string, ingredient *Ingredient) error
	Delete(id string) error
}

type PostgresIngredientRepository struct {
	db *sql.DB
}

func NewPostgresIngredientRepository(db *sql.DB) *PostgresIngredientRepository {
	return &PostgresIngredientRepository{
		db: db,
	}
}

func (r *PostgresIngredientRepository) GetAll() ([]Ingredient, error) {
	rows, err := r.db.Query("SELECT id, name FROM ingredients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := []Ingredient{}
	for rows.Next() {
		var ingredient Ingredient
		if err := rows.Scan(&ingredient.ID, &ingredient.Name); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *PostgresIngredientRepository) GetByID(id string) (*Ingredient, error) {
	var ingredient Ingredient
	err := r.db.QueryRow("SELECT id, name FROM ingredients WHERE id = $1", id).Scan(&ingredient.ID, &ingredient.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &ingredient, nil
}

func (r *PostgresIngredientRepository) Create(ingredient *Ingredient) error {
	err := r.db.QueryRow("INSERT INTO ingredients (name) VALUES ($1) RETURNING id", ingredient.Name).Scan(&ingredient.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresIngredientRepository) Update(id string, ingredient *Ingredient) error {
	_, err := r.db.Exec("UPDATE ingredients SET name = $1 WHERE id = $2", ingredient.Name, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresIngredientRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM ingredients WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
