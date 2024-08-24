package controller

import (
	"errors"
	"fmt"
	"image"
	"io"

	"github.com/harperreed/goflipdot/internal/packet"
	"github.com/harperreed/goflipdot/internal/sign"
)

var (
	ErrSignAlreadyExists = errors.New("sign with this name already exists")
	ErrSignNotFound      = errors.New("sign not found")
	ErrInvalidImage      = errors.New("invalid image for sign")
)

// HanoverController controls one or more Hanover signs
type HanoverController struct {
	port  io.Writer
	signs map[string]*sign.HanoverSign
}

// NewHanoverController creates a new HanoverController
func NewHanoverController(port io.Writer) (*HanoverController, error) {
	if port == nil {
		return nil, errors.New("port cannot be nil")
	}
	return &HanoverController{
		port:  port,
		signs: make(map[string]*sign.HanoverSign),
	}, nil
}

// AddSign adds a sign for the controller to communicate with
func (c *HanoverController) AddSign(name string, sign *sign.HanoverSign) error {
	if sign == nil {
		return errors.New("sign cannot be nil")
	}
	if _, exists := c.signs[name]; exists {
		return fmt.Errorf("%w: %s", ErrSignAlreadyExists, name)
	}
	c.signs[name] = sign
	return nil
}

// StartTestSigns broadcasts the test signs start command
func (c *HanoverController) StartTestSigns() error {
	return c.write(packet.TestSignsStartPacket{})
}

// StopTestSigns broadcasts the test signs stop command
func (c *HanoverController) StopTestSigns() error {
	return c.write(packet.TestSignsStopPacket{})
}

// DrawImage sends an image to a sign to be displayed
func (c *HanoverController) DrawImage(img *image.Gray, signName string) error {
	if img == nil {
		return errors.New("image cannot be nil")
	}
	sign, err := c.getSign(signName)
	if err != nil {
		return err
	}

	if err := sign.ValidateImage(img); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidImage, err)
	}

	flippedImg := sign.FlipImage(img)
	pkt := packet.ImagePacket{
		Address: sign.Address,
		Image:   flippedImg,
	}

	return c.write(pkt)
}

// GetSign returns a sign by name
func (c *HanoverController) GetSign(name string) (*sign.HanoverSign, error) {
	return c.getSign(name)
}

func (c *HanoverController) getSign(name string) (*sign.HanoverSign, error) {
	if name == "" && len(c.signs) == 1 {
		for _, s := range c.signs {
			return s, nil
		}
	}
	if s, ok := c.signs[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrSignNotFound, name)
}

func (c *HanoverController) write(pkt packet.Packet) error {
	bytes, err := pkt.GetBytes()
	if err != nil {
		return fmt.Errorf("failed to get packet bytes: %w", err)
	}
	_, err = c.port.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}
	return nil
}
