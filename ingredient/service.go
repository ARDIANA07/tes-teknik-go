package ingredient

import "github.com/gin-gonic/gin"

type IngredientService struct {
	ingredientRepo IngredientRepository
}

func NewIngredientService(repo IngredientRepository) *IngredientService {
	return &IngredientService{
		ingredientRepo: repo,
	}
}

func (s *IngredientService) GetAllIngredients(c *gin.Context, service *IngredientService) (interface{}, error) {
	ingredients, err := s.ingredientRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (s *IngredientService) CreateIngredient(c *gin.Context, service *IngredientService) (interface{}, error) {
	var ingredient Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		return nil, err
	}

	if err := s.ingredientRepo.Create(&ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (s *IngredientService) GetIngredientByID(c *gin.Context, service *IngredientService) (interface{}, error) {
	id := c.Param("id")

	ingredient, err := s.ingredientRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (s *IngredientService) UpdateIngredient(c *gin.Context, service *IngredientService) (interface{}, error) {
	id := c.Param("id")

	var ingredient Ingredient
	if err := c.ShouldBindJSON(&ingredient); err != nil {
		return nil, err
	}

	if err := s.ingredientRepo.Update(id, &ingredient); err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (s *IngredientService) DeleteIngredient(c *gin.Context, service *IngredientService) (interface{}, error) {
	id := c.Param("id")

	if err := s.ingredientRepo.Delete(id); err != nil {
		return nil, err
	}

	return gin.H{"message": "Ingredient deleted successfully"}, nil
}
