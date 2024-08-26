package packet

import (
	"encoding/binary"
	"errors"
	"image"

	"github.com/harperreed/goflipdot/internal"
)

var (
	ErrInvalidImage = errors.New("invalid image")
)

type Packet interface {
	GetBytes() ([]byte, error)
}

type TestSignsStartPacket struct{}

func (p TestSignsStartPacket) GetBytes() ([]byte, error) {
	return internal.FormatPacket(internal.CommandCodes["start_test_signs"], 0, nil), nil
}

type TestSignsStopPacket struct{}

func (p TestSignsStopPacket) GetBytes() ([]byte, error) {
	return internal.FormatPacket(internal.CommandCodes["stop_test_signs"], 0, nil), nil
}

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
		return nil, err
	}

	payload := make([]byte, 2+len(imageBytes)*2)
	binary.BigEndian.PutUint16(payload[:2], uint16(len(imageBytes)))
	internal.ToAsciiHex(imageBytes)

	return internal.FormatPacket(internal.CommandCodes["write_image"], byte(p.Address), payload), nil
}

func imageToBytes(img *image.Gray) ([]byte, error) {
	if img == nil {
		return nil, ErrInvalidImage
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	byteWidth := (width + 7) / 8
	paddedHeight := internal.ClosestLargerMultiple(height, 8)

	result := make([]byte, byteWidth*paddedHeight/8)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if img.GrayAt(x, y).Y > 127 {
				byteIndex := (paddedHeight - y - 1) / 8 * byteWidth + x/8
				bitIndex := uint(7 - x%8)
				if byteIndex < len(result) {
					result[byteIndex] |= 1 << bitIndex
				}
			}
		}
	}
	return result, nil
}
