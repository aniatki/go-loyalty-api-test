package services

import (
	"github.com/aniatki/loyalty-api/internal/models"
	"github.com/aniatki/loyalty-api/internal/repositories"
)

type ItemService interface {
	CreateItem(input models.CreateItemInput) (*models.Item, error)
	GetItems() ([]models.Item, error)
	UpdateItemTags(itemID uint, tagIDs []uint) (*models.Item, error)
}

type itemService struct {
	itemRepo repositories.ItemRepository
	tagRepo  repositories.TagRepository
}

func NewItemService(itemRepo repositories.ItemRepository, tagRepo repositories.TagRepository) ItemService {
	return &itemService{itemRepo: itemRepo, tagRepo: tagRepo}
}

func (s *itemService) CreateItem(input models.CreateItemInput) (*models.Item, error) {
	item := &models.Item{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Tags:        input.Tags,
	}
	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *itemService) GetItems() ([]models.Item, error) {
	return s.itemRepo.FindAll()
}

func (s *itemService) UpdateItemTags(itemID uint, tagIDs []uint) (*models.Item, error) {
	item, err := s.itemRepo.FindByID(itemID)
	if err != nil {
		return nil, err
	}
	if err := s.itemRepo.UpdateTags(item, tagIDs); err != nil {
		return nil, err
	}
	return item, nil
}
