package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo         *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(repo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

func (s *ProductService) Create(data *models.Product) error {
	if data.CategoryID != 0 {
		categoryExists, err := s.categoryRepo.Exists(data.CategoryID)
		if err != nil {
			return err
		}
		if !categoryExists {
			return errors.New("category not found")
		}
	}

	return s.repo.Create(data)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	_, err := s.repo.GetByID(product.ID)
	if err != nil {
		return errors.New("product not found")
	}
	if product.CategoryID != 0 {
		categoryExists, err := s.categoryRepo.Exists(product.CategoryID)
		if err != nil {
			return err
		}
		if !categoryExists {
			return errors.New("category not found")
		}
	}
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}
