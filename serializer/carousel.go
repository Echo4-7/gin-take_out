package serializer

import "Take_Out/model"

type Carousel struct {
	Id        uint   `json:"id"`
	ImgPath   string `json:"img_path"`
	ProductId uint   `json:"product_id"`
}

func BuildCarousel(item *model.Carousel) {
	
}
