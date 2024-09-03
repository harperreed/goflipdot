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

    bounds := p.Image.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    resolution := (width * height) / 8
    resolutionStr := fmt.Sprintf("%02X", resolution)

    imageBytes, err := imageToBytes(p.Image)
    if err != nil {
        return nil, fmt.Errorf("failed to convert image to bytes: %w", err)
    }

    packet := make([]byte, 0, 5+len(imageBytes)*2+3)
    packet = append(packet, startByte)
    packet = append(packet, '1')
    packet = append(packet, byte(p.Address+'0'))
    packet = append(packet, []byte(resolutionStr)...)

    encodedImageBytes := make([]byte, len(imageBytes)*2)
    hex.Encode(encodedImageBytes, imageBytes)
    packet = append(packet, encodedImageBytes...)

    packet = append(packet, endByte)

    checksum := calculateChecksum(packet)
    checksumStr := fmt.Sprintf("%02X", checksum)
    packet = append(packet, []byte(checksumStr)...)

    return packet, nil
}

func calculateChecksum(data []byte) byte {
    var sum int
    for _, b := range data {
        sum += int(b)
    }
    sum -= int(startByte)
    sum = sum & 0xFF
    return byte((sum ^ 0xFF) + 1)
}

func imageToBytes(img *image.Gray) ([]byte, error) {
    bounds := img.Bounds()
    width, height := bounds.Dx(), bounds.Dy()
    bytesPerColumn := (height + 7) / 8
    result := make([]byte, width*bytesPerColumn)

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            if img.GrayAt(x, height-1-y).Y > 127 { // Flip vertically
                byteIndex := x*bytesPerColumn + (y / 8)
                bitIndex := uint(y % 8)
                result[byteIndex] |= 1 << bitIndex
            }
        }
    }

    return result, nil
}
