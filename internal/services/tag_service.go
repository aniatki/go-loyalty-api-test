package services

import (
	"loyalty-api/internal/models"
	"strings"
)

type TagService interface {
	CreateTag(tag models.Tag) (*models.Tag, error)
	GetTags() ([]models.Tag, error)
	DeleteTag(id uint) error
}

type tagService struct {
	repo repositories.TagRepository
}

func NewTagService(repo repositories.TagRepository) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) formatTagName(name string) string {
	name = strings.ToLower(name)
	words := strings.Fields(name)
	return strings.Join(words, "")
}

func (s *tagService) CreateTag(tag models.Tag) (*models.Tag, error) {
	tag.Name = s.formatTagName(tag.Name)
	existing, err := s.repo.FindByName(tag.Name)
	if existing != nil {
		return nil, fmt.Errorf("tag already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err := s.repo.Create(&tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *tagService) GetTags() ([]models.Tag, error) {
	return s.repo.FindAll()
}

func (s *tagService) DeleteTag(id uint) error {
	return s.repo.Delete(id)
}
