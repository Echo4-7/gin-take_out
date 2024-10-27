package service

//
//import (
//	"Take_Out/dao"
//	"Take_Out/pkg/e"
//	"Take_Out/serializer"
//)
//
//type CarouselService struct {
//}
//
//func (service *CarouselService) List(ctx context.Context) serializer.Response {
//	carouselDao := dao.NewCarouselDao(ctx)
//	code := e.SUCCESS
//	carousels, err := carouselDao.ListCarousel()
//	if err != nil {
//		code = e.ERROR
//		return serializer.Response{
//			Status: code,
//			Msg:    e.GetMsg(code),
//		}
//	}
//	return serializer.BuildListResponse(serializer.BuildCarousel(carousels), uint(len(carousels)))
//}
