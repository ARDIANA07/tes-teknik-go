package main

import (
	"app_restoran/category"
	"app_restoran/ingredient"
	"app_restoran/recipe"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", "postgres://postgres:ardi1234@localhost:5432/db_restaurant?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create ingredient repository and service
	ingredientRepo := ingredient.NewPostgresIngredientRepository(db)
	ingredientService := ingredient.NewIngredientService(ingredientRepo)

	// Create category repository and service
	categoryRepo := category.NewPostgresCategoryRepository(db)
	categoryService := category.NewCategoryService(categoryRepo)

	// Create recipe repository and service
	recipeRepo := recipe.NewPostgresRecipeRepository(db)
	recipeService := recipe.NewRecipeService(recipeRepo)

	router := gin.Default()

	// Ingredient routes
	router.GET("/ingredients", ingredientHandler(ingredientService.GetAllIngredients))
	router.POST("/ingredients", ingredientHandler(ingredientService.CreateIngredient))
	router.GET("/ingredients/:id", ingredientHandler(ingredientService.GetIngredientByID))
	router.PUT("/ingredients/:id", ingredientHandler(ingredientService.UpdateIngredient))
	router.DELETE("/ingredients/:id", ingredientHandler(ingredientService.DeleteIngredient))

	// Category routes
	router.GET("/categories", categoryHandler(categoryService.GetAllCategories))
	router.POST("/categories", categoryHandler(categoryService.CreateCategory))
	router.GET("/categories/:id", categoryHandler(categoryService.GetCategoryByID))
	router.PUT("/categories/:id", categoryHandler(categoryService.UpdateCategory))
	router.DELETE("/categories/:id", categoryHandler(categoryService.DeleteCategory))

	// Recipe routes
	router.GET("/recipes", recipeHandler(recipeService.GetAllRecipes))
	router.POST("/recipes", recipeHandler(recipeService.CreateRecipe))
	router.GET("/recipes/:id", recipeHandler(recipeService.GetRecipeByID))
	router.PUT("/recipes/:id", recipeHandler(recipeService.UpdateRecipe))
	router.DELETE("/recipes/:id", recipeHandler(recipeService.DeleteRecipe))
	router.GET("/recipes", recipeHandler(recipeService.GetAllRecipes))
	router.GET("/recipes/:id", recipeHandler(recipeService.GetRecipeByID))
	router.GET("/recipes/search", recipeHandler(recipeService.SearchRecipes))

	router.Run(":8080")
}

type ingredientHandlerFunc func(c *gin.Context, service *ingredient.IngredientService) (interface{}, error)

func ingredientHandler(handler ingredientHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := ingredient.NewIngredientService(ingredient.NewPostgresIngredientRepository(db))

		data, err := handler(c, service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

type categoryHandlerFunc func(c *gin.Context, service *category.CategoryService) (interface{}, error)

func categoryHandler(handler categoryHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := category.NewCategoryService(category.NewPostgresCategoryRepository(db))

		data, err := handler(c, service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}

type recipeHandlerFunc func(c *gin.Context, service *recipe.RecipeService) (interface{}, error)

func recipeHandler(handler recipeHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := recipe.NewRecipeService(recipe.NewPostgresRecipeRepository(db))

		data, err := handler(c, service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, data)
	}
}
