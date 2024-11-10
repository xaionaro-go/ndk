package main

import (
	"fmt"
	"image"

	"github.com/xaionaro-go/camera/ximage"
	media "github.com/xaionaro-go/ndk/media24"
)

func ConvertImage(
	img *media.Image,
) (image.Image, error) {
	//defer img.Delete()

	w, _ := img.GetWidth()
	h, _ := img.GetHeight()
	//t, _ := img.GetTimestamp()
	format, _ := img.GetFormat()
	//planes, _ := img.GetNumberOfPlanes()
	//pixelStride, _ := img.GetPlanePixelStride(1)
	//rowStride, _ := img.GetPlaneRowStride(1)
	switch format {
	case media.FORMAT_YUV_420_888:
		dataY, err := img.GetPlaneData(0)
		if err != nil {
			return nil, fmt.Errorf("img.GetPlaneData(0): %w", err)
		}
		if len(dataY) == 0 {
			return nil, fmt.Errorf("len(dataY) == 0")
		}
		dataUV, err := img.GetPlaneData(2)
		if err != nil {
			return nil, fmt.Errorf("img.GetPlaneData(2): %w", err)
		}
		if len(dataY) == 0 {
			return nil, fmt.Errorf("len(dataUV) == 0")
		}
		img := ximage.NewNV12NoAlloc(image.Rectangle{Max: image.Point{
			X: w,
			Y: h,
		}})
		img.Y = dataY
		if len(dataUV)%2 == 1 { // TODO: investigate why sometimes the length is one byte short
			dataUV = append(dataUV, 0)
		}
		if err := img.SetCbCrBytes(dataUV); err != nil {
			return nil, fmt.Errorf("SetCbCrBytes: %w", err)
		}
		return img, nil
	default:
		return nil, fmt.Errorf("not supported: %v", format)
	}
}
