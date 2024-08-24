package test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/harperreed/goflipdot/internal/packet"
)

func TestPackets(t *testing.T) {
	t.Run("TestSignsStartPacket", func(t *testing.T) {
		p := packet.TestSignsStartPacket{}
		gotBytes, err := p.GetBytes()
		if err != nil {
			t.Fatalf("Failed to get bytes: %v", err)
		}
		expected := []byte{0x02, '3', '0', 0x03, '9', 'A'}
		if !bytes.Equal(gotBytes, expected) {
			t.Errorf("Unexpected TestSignsStartPacket bytes. Got %v, want %v", gotBytes, expected)
		}
	})

	t.Run("TestSignsStopPacket", func(t *testing.T) {
		p := packet.TestSignsStopPacket{}
		gotBytes, err := p.GetBytes()
		if err != nil {
			t.Fatalf("Failed to get bytes: %v", err)
		}
		expected := []byte{0x02, 'C', '0', 0x03, '8', 'A'}
		if !bytes.Equal(gotBytes, expected) {
			t.Errorf("Unexpected TestSignsStopPacket bytes. Got %v, want %v", gotBytes, expected)
		}
	})

	t.Run("ImagePacket", func(t *testing.T) {
		img := image.NewGray(image.Rect(0, 0, 8, 8))
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				if (x+y)%2 == 0 {
					img.Set(x, y, color.White)
				}
			}
		}

		p := packet.ImagePacket{
			Address: 1,
			Image:   img,
		}

		gotBytes, err := p.GetBytes()
		if err != nil {
			t.Fatalf("Failed to get bytes: %v", err)
		}

		expectedLength := 1 + 2 + 2 + 16 + 1 + 2 // Start byte + address + data length + image data + end byte + checksum
		if len(gotBytes) != expectedLength {
			t.Errorf("Unexpected ImagePacket length. Got %d, want %d", len(gotBytes), expectedLength)
		}

		if gotBytes[0] != 0x02 || gotBytes[1] != '1' || gotBytes[2] != '1' {
			t.Errorf("Unexpected ImagePacket header. Got %v", gotBytes[:3])
		}

		if gotBytes[len(gotBytes)-3] != 0x03 {
			t.Errorf("Unexpected ImagePacket end byte. Got %v", gotBytes[len(gotBytes)-3])
		}
	})
}
