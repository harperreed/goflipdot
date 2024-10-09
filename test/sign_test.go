package test

import (
	"image"
	"image/color"
	"testing"

	"github.com/harperreed/goflipdot/internal/sign"
)

func TestNewHanoverSign(t *testing.T) {
	testCases := []struct {
		name        string
		address     int
		width       int
		height      int
		flip        bool
		expectError bool
	}{
		{"ValidSign", 1, 86, 7, false, false},
		{"InvalidAddress", -1, 86, 7, false, true},
		{"InvalidWidth", 1, 0, 7, false, true},
		{"InvalidHeight", 1, 86, 0, false, true},
		{"LargeSign", 2, 200, 50, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := sign.NewHanoverSign(tc.address, tc.width, tc.height, tc.flip)
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if s.Address != tc.address || s.Width != tc.width || s.Height != tc.height || s.Flip != tc.flip {
					t.Errorf("Sign properties do not match. Got %+v, want %+v", s, tc)
				}
			}
		})
	}
}

func TestCreateImage(t *testing.T) {
	testCases := []struct {
		name   string
		width  int
		height int
	}{
		{"SmallSign", 10, 5},
		{"StandardSign", 86, 7},
		{"LargeSign", 200, 50},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := sign.NewHanoverSign(1, tc.width, tc.height, false)
			img := s.CreateImage()
			if img.Bounds().Dx() != tc.width || img.Bounds().Dy() != tc.height {
				t.Errorf("Unexpected image size. Got %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), tc.width, tc.height)
			}
		})
	}
}

func TestValidateImage(t *testing.T) {
	s, _ := sign.NewHanoverSign(1, 86, 7, false)

	testCases := []struct {
		name        string
		img         *image.Gray
		expectError bool
	}{
		{"ValidImage", image.NewGray(image.Rect(0, 0, 86, 7)), false},
		{"NilImage", nil, true},
		{"WrongWidth", image.NewGray(image.Rect(0, 0, 100, 7)), true},
		{"WrongHeight", image.NewGray(image.Rect(0, 0, 86, 10)), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.ValidateImage(tc.img)
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFlipImage(t *testing.T) {
	testCases := []struct {
		name         string
		flip         bool
		setPixel     func(*image.Gray)
		checkFlipped func(*testing.T, *image.Gray)
	}{
		{
			name: "NoFlip",
			flip: false,
			setPixel: func(img *image.Gray) {
				img.Set(0, 0, color.White)
				img.Set(85, 6, color.Black)
			},
			checkFlipped: func(t *testing.T, img *image.Gray) {
				if img.At(0, 0).(color.Gray).Y != 255 {
					t.Error("Bottom-left pixel should be white")
				}
				if img.At(85, 6).(color.Gray).Y != 0 {
					t.Error("Top-right pixel should be black")
				}
			},
		},
		{
			name: "Flip",
			flip: true,
			setPixel: func(img *image.Gray) {
				img.Set(0, 0, color.White)
				img.Set(85, 6, color.Black)
			},
			checkFlipped: func(t *testing.T, img *image.Gray) {
				if img.At(0, 6).(color.Gray).Y != 255 {
					t.Error("Top-left pixel should be white after flipping")
				}
				if img.At(85, 0).(color.Gray).Y != 0 {
					t.Error("Bottom-right pixel should be black after flipping")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, _ := sign.NewHanoverSign(1, 86, 7, tc.flip)
			img := image.NewGray(image.Rect(0, 0, 86, 7))
			tc.setPixel(img)
			flippedImg := s.FlipImage(img)
			tc.checkFlipped(t, flippedImg)
		})
	}
}
