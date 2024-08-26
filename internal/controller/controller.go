package controller

import (
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"time"

	"github.com/tarm/serial"
	"github.com/harperreed/goflipdot/internal"
	"github.com/harperreed/goflipdot/internal/packet"
	"github.com/harperreed/goflipdot/internal/sign"
)

var (
	ErrSignAlreadyExists = errors.New("sign with this name already exists")
	ErrSignNotFound      = errors.New("sign not found")
	ErrInvalidImage      = errors.New("invalid image for sign")
)

type HanoverController struct {
	port  *serial.Port
	signs map[string]*sign.HanoverSign
}

func NewHanoverController(portName string) (*HanoverController, error) {
	config := &serial.Config{
		Name:        portName,
		Baud:        internal.BaudRate,
		ReadTimeout: time.Second * 5, // 5 second timeout, adjust as needed
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}
	return &HanoverController{
		port:  port,
		signs: make(map[string]*sign.HanoverSign),
	}, nil
}


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

func (c *HanoverController) StartTestSigns() error {
	return c.write(packet.TestSignsStartPacket{})
}

func (c *HanoverController) StopTestSigns() error {
	return c.write(packet.TestSignsStopPacket{})
}

func (c *HanoverController) DrawImage(img *image.Gray, signName string) error {
	if img == nil {
		return errors.New("image cannot be nil")
	}
	sign, err := c.getSign(signName)
	if err != nil {
		return err
	}

	log.Printf("Drawing image for sign %s: %dx%d", signName, img.Bounds().Dx(), img.Bounds().Dy())

	pkt, err := sign.ToImagePacket(img)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidImage, err)
	}

	return c.write(pkt)
}

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
	log.Printf("Writing packet: %v", bytes)
	_, err = c.port.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}
	return nil
}
