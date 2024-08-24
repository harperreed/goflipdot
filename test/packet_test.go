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
		expected := []byte{0x02, '3', '0', 0x03, '9', 'A'}
		if !bytes.Equal(p.GetBytes(), expected) {
			t.Errorf("Unexpected TestSignsStartPacket bytes. Got %v, want %v", p.GetBytes(), expected)
		}
	})

	t.Run("TestSignsStopPacket", func(t *testing.T) {
		p := packet.TestSignsStopPacket{}
		expected := []byte{0x02, 'C', '0', 0x03, '8', 'A'}
		if !bytes.Equal(p.GetBytes(), expected) {
			t.Errorf("Unexpected TestSignsStopPacket bytes. Got %v, want %v", p.GetBytes(), expected)
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

		bytes := p.GetBytes()
		if len(bytes) != 25 {
			t.Errorf("Unexpected ImagePacket length. Got %d, want 25", len(bytes))
		}

		if bytes[0] != 0x02 || bytes[1] != '1' || bytes[2] != '1' {
			t.Errorf("Unexpected ImagePacket header. Got %v", bytes[:3])
		}

		if bytes[len(bytes)-3] != 0x03 {
			t.Errorf("Unexpected ImagePacket end byte. Got %v", bytes[len(bytes)-3])
		}
	})
}
