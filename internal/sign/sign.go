package sign

import (
	"errors"
	"image"
)

// HanoverSign represents a Hanover flipdot sign
type HanoverSign struct {
	Address int
	Width   int
	Height  int
	Flip    bool
}

// NewHanoverSign creates a new HanoverSign
func NewHanoverSign(address, width, height int, flip bool) *HanoverSign {
	return &HanoverSign{
		Address: address,
		Width:   width,
		Height:  height,
		Flip:    flip,
	}
}

// CreateImage creates a blank image for the sign
func (s *HanoverSign) CreateImage() *image.Gray {
	return image.NewGray(image.Rect(0, 0, s.Width, s.Height))
}

// ValidateImage checks if the given image is compatible with the sign
func (s *HanoverSign) ValidateImage(img *image.Gray) error {
	bounds := img.Bounds()
	if bounds.Dx() != s.Width || bounds.Dy() != s.Height {
		return errors.New("image dimensions do not match sign dimensions")
	}
	return nil
}

// FlipImage rotates the image 180 degrees if the sign is flipped
func (s *HanoverSign) FlipImage(img *image.Gray) *image.Gray {
	if !s.Flip {
		return img
	}

	flipped := image.NewGray(img.Bounds())
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			flipped.Set(x, y, img.At(s.Width-1-x, s.Height-1-y))
		}
	}
	return flipped
}
