package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Tags        []Tag   `json:"tags" gorm:"many2many:item_tags"`
}

type Tag struct {
	gorm.Model
	Name string `json:"name" gorm:"unique; not null"`
}

type CreateItemInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Tags        []Tag   `json:"tags"`
}

type UpdateItemTagsInput struct {
	TagIDs []uint `json:"tag_ids"`
}
