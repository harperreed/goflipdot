package sign

import (
	"errors"
	"image"
)

type HanoverSign struct {
	Address int
	Width   int
	Height  int
}

func NewHanoverSign(address, width, height int) (*HanoverSign, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width and height must be positive")
	}
	if address < 0 {
		return nil, errors.New("address must be non-negative")
	}
	return &HanoverSign{
		Address: address,
		Width:   width,
		Height:  height,
	}, nil
}

func (s *HanoverSign) CreateImage() *image.Gray {
	return image.NewGray(image.Rect(0, 0, s.Width, s.Height))
}

func (s *HanoverSign) ValidateImage(img *image.Gray) error {
	if img == nil {
		return errors.New("image cannot be nil")
	}
	bounds := img.Bounds()
	if bounds.Dx() != s.Width || bounds.Dy() != s.Height {
		return errors.New("image dimensions do not match sign dimensions")
	}
	return nil
}
