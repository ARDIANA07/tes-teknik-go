package recipe

import "github.com/gin-gonic/gin"

type RecipeService struct {
	recipeRepo RecipeRepository
}

func NewRecipeService(repo RecipeRepository) *RecipeService {
	return &RecipeService{
		recipeRepo: repo,
	}
}

func (s *RecipeService) GetAllRecipes(c *gin.Context, service *RecipeService) (interface{}, error) {
	recipes, err := s.recipeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (s *RecipeService) CreateRecipe(c *gin.Context, service *RecipeService) (interface{}, error) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		return nil, err
	}

	if err := s.recipeRepo.Create(&recipe); err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *RecipeService) GetRecipeByID(c *gin.Context, service *RecipeService) (interface{}, error) {
	id := c.Param("id")

	recipe, err := s.recipeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *RecipeService) UpdateRecipe(c *gin.Context, service *RecipeService) (interface{}, error) {
	id := c.Param("id")

	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		return nil, err
	}

	if err := s.recipeRepo.Update(id, &recipe); err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *RecipeService) DeleteRecipe(c *gin.Context, service *RecipeService) (interface{}, error) {
	id := c.Param("id")

	if err := s.recipeRepo.Delete(id); err != nil {
		return nil, err
	}

	return gin.H{"message": "Recipe deleted successfully"}, nil
}

func (s *RecipeService) SearchRecipes(c *gin.Context, service *RecipeService) (interface{}, error) {
	ingredients := c.QueryArray("ingredients")
	category := c.Query("category")

	recipes, err := s.recipeRepo.SearchRecipes(ingredients, category)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}
