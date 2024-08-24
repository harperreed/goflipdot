package controller

import (
	"errors"
	"image"
	"io"

	"github.com/harperreed/goflipdot/internal/packet"
	"github.com/harperreed/goflipdot/internal/sign"
)

// HanoverController controls one or more Hanover signs
type HanoverController struct {
	port  io.Writer
	signs map[string]*sign.HanoverSign
}

// NewHanoverController creates a new HanoverController
func NewHanoverController(port io.Writer) *HanoverController {
	return &HanoverController{
		port:  port,
		signs: make(map[string]*sign.HanoverSign),
	}
}

// AddSign adds a sign for the controller to communicate with
func (c *HanoverController) AddSign(name string, sign *sign.HanoverSign) error {
	if _, exists := c.signs[name]; exists {
		return errors.New("sign with this name already exists")
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
	sign, err := c.getSign(signName)
	if err != nil {
		return err
	}

	if err := sign.ValidateImage(img); err != nil {
		return err
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
	return nil, errors.New("sign not found")
}

func (c *HanoverController) write(pkt packet.Packet) error {
	_, err := c.port.Write(pkt.GetBytes())
	return err
}
