package packet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image"
)

const (
	startByte byte = 0x02
	endByte   byte = 0x03
)

var (
	ErrInvalidImage = errors.New("invalid image")
)

// Packet represents a data packet for Hanover signs
type Packet interface {
	GetBytes() ([]byte, error)
}

// TestSignsStartPacket is a command for all signs to cycle through a test mode sequence
type TestSignsStartPacket struct{}

func (p TestSignsStartPacket) GetBytes() ([]byte, error) {
	return []byte{startByte, '3', '0', endByte, '9', 'A'}, nil
}

// TestSignsStopPacket is a command for all signs to stop test mode sequence
type TestSignsStopPacket struct{}

func (p TestSignsStopPacket) GetBytes() ([]byte, error) {
	return []byte{startByte, 'C', '0', endByte, '8', 'A'}, nil
}

// ImagePacket encodes an image to display
type ImagePacket struct {
	Address int
	Image   *image.Gray
}

func (p ImagePacket) GetBytes() ([]byte, error) {
	if p.Image == nil {
		return nil, ErrInvalidImage
	}

	imageBytes, err := imageToBytes(p.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to bytes: %w", err)
	}

	payload := make([]byte, 3+len(imageBytes)*2)
	payload[0] = byte(p.Image.Bounds().Dx() & 0xFF)
	payload[1] = byte(p.Image.Bounds().Dx() >> 8)
	payload[2] = byte(p.Image.Bounds().Dy() & 0xFF)
	payload[3] = byte(p.Image.Bounds().Dy() >> 8)

	width := fmt.Sprintf("%04d", p.Image.Bounds().Dx())
	height := fmt.Sprintf("%04d", p.Image.Bounds().Dy())
	copy(payload[0:4], width)
	copy(payload[4:8], height)

	hex.Encode(payload[3:], imageBytes)

	packet := make([]byte, 0, 5+len(payload))
	packet = append(packet, startByte)
	packet = append(packet, []byte(fmt.Sprintf("1%X", p.Address))...)
	packet = append(packet, payload...)
	packet = append(packet, endByte)

	checksum := calculateChecksum(packet)
	packet = append(packet, []byte(fmt.Sprintf("%02X", checksum))...)

	return packet, nil
}

func calculateChecksum(data []byte) byte {
	var sum int
	for _, b := range data[1 : len(data)-1] { // Exclude start byte and end byte
		sum += int(b)
	}
	sum += int(endByte) // Add end byte to sum
	sum = sum & 0xFF    // Truncate to 8 bits
	return byte((sum ^ 0xFF) + 1)
}

func imageToBytes(img *image.Gray) ([]byte, error) {
	if img == nil {
		return nil, ErrInvalidImage
	}

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

	// Flip the byte order within each column
	for col := 0; col < byteWidth; col++ {
		for i := 0; i < height/2; i++ {
			result[col*height+i], result[col*height+(height-1-i)] = result[col*height+(height-1-i)], result[col*height+i]
		}
	}

	return result, nil
}
