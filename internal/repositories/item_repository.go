package repositories

import (
	"github.com/aniatki/loyalty-api/internal/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(item *models.Item) error
	FindAll() ([]models.Item, error)
	FindByID(id uint) (*models.Item, error)
	UpdateTags(item *models.Item, tagIDs []uint) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Preload("Tags").Find(&items).Error
	return items, err
}

func (r *itemRepository) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.Preload("Tags").First(&item, id).Error
	return &item, err
}

func (r *itemRepository) UpdateTags(item *models.Item, tagIDs []uint) error {
	var tags []models.Tag
	if len(tagIDs) > 0 {
		if err := r.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
	}
	return r.db.Model(item).Association("Tags").Replace(tags)
}
