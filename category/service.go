package category

import "github.com/gin-gonic/gin"

type CategoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: repo,
	}
}

func (s *CategoryService) GetAllCategories(c *gin.Context, service *CategoryService) (interface{}, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) CreateCategory(c *gin.Context, service *CategoryService) (interface{}, error) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		return nil, err
	}

	if err := s.categoryRepo.Create(&category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetCategoryByID(c *gin.Context, service *CategoryService) (interface{}, error) {
	id := c.Param("id")

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(c *gin.Context, service *CategoryService) (interface{}, error) {
	id := c.Param("id")

	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		return nil, err
	}

	if err := s.categoryRepo.Update(id, &category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(c *gin.Context, service *CategoryService) (interface{}, error) {
	id := c.Param("id")

	if err := s.categoryRepo.Delete(id); err != nil {
		return nil, err
	}

	return gin.H{"message": "Category deleted successfully"}, nil
}
