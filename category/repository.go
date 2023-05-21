package category

import "database/sql"

type CategoryRepository interface {
	GetAll() ([]Category, error)
	GetByID(id string) (*Category, error)
	Create(category *Category) error
	Update(id string, category *Category) error
	Delete(id string) error
}

type PostgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{
		db: db,
	}
}

func (r *PostgresCategoryRepository) GetAll() ([]Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *PostgresCategoryRepository) GetByID(id string) (*Category, error) {
	var category Category
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &category, nil
}

func (r *PostgresCategoryRepository) Create(category *Category) error {
	err := r.db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", category.Name).Scan(&category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresCategoryRepository) Update(id string, category *Category) error {
	_, err := r.db.Exec("UPDATE categories SET name = $1 WHERE id = $2", category.Name, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresCategoryRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
