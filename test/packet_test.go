package test

import (
	"bytes"
	"image"
	"image/color"
	"math/rand"
	"testing"

	"github.com/harperreed/goflipdot/internal/packet"
)

func TestTestSignsStartPacket(t *testing.T) {
	p := packet.TestSignsStartPacket{}
	gotBytes, err := p.GetBytes()
	if err != nil {
		t.Fatalf("Failed to get bytes: %v", err)
	}
	expected := []byte{0x02, '3', '0', 0x03, '9', 'A'}
	if !bytes.Equal(gotBytes, expected) {
		t.Errorf("Unexpected TestSignsStartPacket bytes. Got %v, want %v", gotBytes, expected)
	}
}

func TestTestSignsStopPacket(t *testing.T) {
	p := packet.TestSignsStopPacket{}
	gotBytes, err := p.GetBytes()
	if err != nil {
		t.Fatalf("Failed to get bytes: %v", err)
	}
	expected := []byte{0x02, 'C', '0', 0x03, '8', 'A'}
	if !bytes.Equal(gotBytes, expected) {
		t.Errorf("Unexpected TestSignsStopPacket bytes. Got %v, want %v", gotBytes, expected)
	}
}

func TestImagePacket(t *testing.T) {
	testCases := []struct {
		name        string
		width       int
		height      int
		pattern     func(image.Image)
		expectError bool
	}{
		{"SmallCheckerboard", 8, 8, checkerboardPattern, false},
		{"LargeRandom", 86, 7, randomPattern, false},
		{"InvalidSize", 0, 0, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := image.NewGray(image.Rect(0, 0, tc.width, tc.height))
			if tc.pattern != nil {
				tc.pattern(img)
			}

			p := packet.ImagePacket{
				Address: 1,
				Image:   img,
			}

			gotBytes, err := p.GetBytes()

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to get bytes: %v", err)
			}

			expectedLength := 1 + 2 + 2 + (tc.width*tc.height+7)/8 + 1 + 2
			if len(gotBytes) != expectedLength {
				t.Errorf("Unexpected ImagePacket length. Got %d, want %d", len(gotBytes), expectedLength)
			}

			if gotBytes[0] != 0x02 || gotBytes[1] != '1' || gotBytes[2] != '1' {
				t.Errorf("Unexpected ImagePacket header. Got %v", gotBytes[:3])
			}

			if gotBytes[len(gotBytes)-3] != 0x03 {
				t.Errorf("Unexpected ImagePacket end byte. Got %v", gotBytes[len(gotBytes)-3])
			}

			// Verify checksum
			checksum := calculateChecksum(gotBytes[:len(gotBytes)-2])
			if gotBytes[len(gotBytes)-2] != checksum[0] || gotBytes[len(gotBytes)-1] != checksum[1] {
				t.Errorf("Invalid checksum. Got %v, want %v", gotBytes[len(gotBytes)-2:], checksum)
			}
		})
	}
}

func checkerboardPattern(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if (x+y)%2 == 0 {
				img.(*image.Gray).Set(x, y, color.White)
			}
		}
	}
}

func randomPattern(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if rand.Intn(2) == 0 {
				img.(*image.Gray).Set(x, y, color.White)
			}
		}
	}
}

func calculateChecksum(data []byte) []byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	sum = (^sum + 1) & 0xFF
	return []byte{sum / 16, sum % 16}
}
