package recipe

import (
	"database/sql"
	"strconv"
)

type RecipeRepository interface {
	GetAll() ([]Recipe, error)
	GetByID(id string) (*Recipe, error)
	Create(recipe *Recipe) error
	Update(id string, recipe *Recipe) error
	Delete(id string) error

	GetAllRecipes() ([]Recipe, error)
	GetRecipeByID(id int) (*Recipe, error)
	SearchRecipes(ingredients []string, category string) ([]Recipe, error)
}

type PostgresRecipeRepository struct {
	db *sql.DB
}

func NewPostgresRecipeRepository(db *sql.DB) *PostgresRecipeRepository {
	return &PostgresRecipeRepository{
		db: db,
	}
}

func (r *PostgresRecipeRepository) GetAll() ([]Recipe, error) {
	rows, err := r.db.Query("SELECT id, name FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipes := []Recipe{}
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(&recipe.ID, &recipe.Name); err != nil {
			return nil, err
		}

		// Get ingredients for the recipe
		recipe.Ingredients, err = r.getIngredientsByRecipeID(recipe.ID)
		if err != nil {
			return nil, err
		}

		// Get category for the recipe
		recipe.Category, err = r.getCategoryByRecipeID(recipe.ID)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *PostgresRecipeRepository) GetByID(id string) (*Recipe, error) {
	var recipe Recipe
	err := r.db.QueryRow("SELECT id, name FROM recipes WHERE id = $1", id).Scan(&recipe.ID, &recipe.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	// Get ingredients for the recipe
	recipe.Ingredients, err = r.getIngredientsByRecipeID(recipe.ID)
	if err != nil {
		return nil, err
	}

	// Get category for the recipe
	recipe.Category, err = r.getCategoryByRecipeID(recipe.ID)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (r *PostgresRecipeRepository) Create(recipe *Recipe) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow("INSERT INTO recipes (name) VALUES ($1) RETURNING id", recipe.Name).Scan(&recipe.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		_, err := tx.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient) VALUES ($1, $2)", recipe.ID, ingredient)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *PostgresRecipeRepository) Update(id string, recipe *Recipe) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE recipes SET name = $1 WHERE id = $2", recipe.Name, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		_, err := tx.Exec("INSERT INTO recipe_ingredients (recipe_id, ingredient) VALUES ($1, $2)", id, ingredient)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *PostgresRecipeRepository) Delete(id string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM recipes WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM recipe_ingredients WHERE recipe_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *PostgresRecipeRepository) getIngredientsByRecipeID(recipeID string) ([]string, error) {
	rows, err := r.db.Query("SELECT ingredient FROM recipe_ingredients WHERE recipe_id = $1", recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := []string{}
	for rows.Next() {
		var ingredient string
		if err := rows.Scan(&ingredient); err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, nil
}

func (r *PostgresRecipeRepository) getCategoryByRecipeID(recipeID string) (string, error) {
	var category string
	err := r.db.QueryRow("SELECT category FROM recipe_categories WHERE recipe_id = $1", recipeID).Scan(&category)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return category, nil
}

func (r *PostgresRecipeRepository) GetAllRecipes() ([]Recipe, error) {
	recipes := []Recipe{}

	rows, err := r.db.Query("SELECT id, name FROM recipes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.Name)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (r *PostgresRecipeRepository) GetRecipeByID(id int) (*Recipe, error) {
	recipe := Recipe{}

	row := r.db.QueryRow("SELECT id, name FROM recipes WHERE id = $1", id)
	err := row.Scan(&recipe.ID, &recipe.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &recipe, nil
}

func (r *PostgresRecipeRepository) SearchRecipes(ingredients []string, category string) ([]Recipe, error) {
	recipes := []Recipe{}

	// Prepare the query
	query := "SELECT id, name FROM recipes WHERE 1=1"
	args := []interface{}{}

	// Add filter conditions based on the provided ingredients
	if len(ingredients) > 0 {
		query += " AND id IN (SELECT recipe_id FROM recipe_ingredients WHERE ingredient IN ("
		for i, ingredient := range ingredients {
			query += "$" + strconv.Itoa(i+1) + ","
			args = append(args, ingredient)
		}
		query = query[:len(query)-1] + "))"
	}

	// Add filter condition based on the provided category
	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe Recipe
		err := rows.Scan(&recipe.ID, &recipe.Name)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
