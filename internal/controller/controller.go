package packet

import (
	"encoding/hex"
	"fmt"
	"image"
	"github.com/harperreed/goflipdot/internal/packet"
	"github.com/harperreed/goflipdot/internal/sign"
)

const (
	startByte byte = 0x02
	endByte   byte = 0x03
)

// Packet represents a data packet for Hanover signs
type Packet interface {
	GetBytes() []byte
}

// TestSignsStartPacket is a command for all signs to cycle through a test mode sequence
type TestSignsStartPacket struct{}

func (p TestSignsStartPacket) GetBytes() []byte {
	return []byte{startByte, '3', '0', endByte, '9', 'A'}
}

// TestSignsStopPacket is a command for all signs to stop test mode sequence
type TestSignsStopPacket struct{}

func (p TestSignsStopPacket) GetBytes() []byte {
	return []byte{startByte, 'C', '0', endByte, '8', 'A'}
}

// ImagePacket encodes an image to display
type ImagePacket struct {
	Address int
	Image   *image.Gray
}

func (p ImagePacket) GetBytes() []byte {
	imageBytes := imageToBytes(p.Image)
	payload := make([]byte, 2+len(imageBytes)*2)
	payload[0] = byte(len(imageBytes) & 0xFF)
	payload[1] = byte(len(imageBytes) >> 8)
	hex.Encode(payload[2:], imageBytes)

	packet := make([]byte, 0, 5+len(payload))
	packet = append(packet, startByte)
	packet = append(packet, []byte(fmt.Sprintf("1%X", p.Address))...)
	packet = append(packet, payload...)
	packet = append(packet, endByte)

	checksum := calculateChecksum(packet)
	packet = append(packet, []byte(fmt.Sprintf("%02X", checksum))...)

	return packet
}

func calculateChecksum(data []byte) byte {
	var sum byte
	for _, b := range data[1:] {
		sum += b
	}
	return (^sum + 1) & 0xFF
}

func imageToBytes(img *image.Gray) []byte {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	byteWidth := (width + 7) / 8

	result := make([]byte, byteWidth*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if img.GrayAt(x, y).Y > 127 {
				byteIndex := y*byteWidth + x/8
				bitIndex := uint(7 - x%8)
				result[byteIndex] |= 1 << bitIndex
			}
		}
	}
	return result
}
