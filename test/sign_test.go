package test

import (
	"image"
	"image/color"
	"testing"

	"github.com/harperreed/goflipdot/internal/sign"
)

func TestSign(t *testing.T) {
	s := sign.NewHanoverSign(1, 86, 7, false)

	t.Run("CreateImage", func(t *testing.T) {
		img := s.CreateImage()
		if img.Bounds().Dx() != 86 || img.Bounds().Dy() != 7 {
			t.Errorf("Unexpected image size. Got %dx%d, want 86x7", img.Bounds().Dx(), img.Bounds().Dy())
		}
	})

	t.Run("ValidateImage", func(t *testing.T) {
		goodImg := image.NewGray(image.Rect(0, 0, 86, 7))
		if err := s.ValidateImage(goodImg); err != nil {
			t.Errorf("Unexpected error for valid image: %v", err)
		}

		badImg := image.NewGray(image.Rect(0, 0, 100, 10))
		if err := s.ValidateImage(badImg); err == nil {
			t.Error("Expected error for invalid image, got nil")
		}
	})

	t.Run("FlipImage", func(t *testing.T) {
		s := sign.NewHanoverSign(1, 86, 7, false)
		img := image.NewGray(image.Rect(0, 0, 86, 7))
		img.Set(0, 0, color.White)
		img.Set(85, 6, color.Black)

		flippedImg := s.FlipImage(img)
		if flippedImg.At(0, 0) != img.At(0, 0) {
			t.Error("Image should not be flipped when sign.Flip is false")
		}

		s.Flip = true
		flippedImg = s.FlipImage(img)
		if flippedImg.At(85, 6) != color.White {
			t.Error("Top-right pixel should be white after flipping")
		}
		if flippedImg.At(0, 0) != color.Black {
			t.Error("Bottom-left pixel should be black after flipping")
		}
	})
}
