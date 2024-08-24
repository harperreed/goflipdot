package sign

import (
	"errors"
	"fmt"
	"image"
)

var (
	ErrInvalidDimensions = errors.New("image dimensions do not match sign dimensions")
)

// HanoverSign represents a Hanover flipdot sign
type HanoverSign struct {
	Address int
	Width   int
	Height  int
	Flip    bool
}

// NewHanoverSign creates a new HanoverSign
func NewHanoverSign(address, width, height int, flip bool) (*HanoverSign, error) {
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
		Flip:    flip,
	}, nil
}

// CreateImage creates a blank image for the sign
func (s *HanoverSign) CreateImage() *image.Gray {
	return image.NewGray(image.Rect(0, 0, s.Width, s.Height))
}

// ValidateImage checks if the given image is compatible with the sign
func (s *HanoverSign) ValidateImage(img *image.Gray) error {
	if img == nil {
		return errors.New("image cannot be nil")
	}
	bounds := img.Bounds()
	if bounds.Dx() != s.Width || bounds.Dy() != s.Height {
		return fmt.Errorf("%w: expected %dx%d, got %dx%d",
			ErrInvalidDimensions, s.Width, s.Height, bounds.Dx(), bounds.Dy())
	}
	return nil
}

// FlipImage rotates the image 180 degrees if the sign is flipped
func (s *HanoverSign) FlipImage(img *image.Gray) *image.Gray {
	if img == nil {
		return nil
	}
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
