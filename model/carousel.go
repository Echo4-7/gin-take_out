package model

import "gorm.io/gorm"

type Carousel struct {
	gorm.Model
	picture string
}
