package controller

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"time"
	"github.com/tarm/serial"
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
	port  io.ReadWriter
	signs map[string]*sign.HanoverSign
}

// NewHanoverController creates a new HanoverController
func NewHanoverController(serialPort string) (*HanoverController, error) {
	c := &serial.Config{Name: serialPort, Baud: 4800}
	port, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}
	return &HanoverController{
		port:  port,
		signs: make(map[string]*sign.HanoverSign),
	}, nil
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
	return c.writeAndRead(packet.TestSignsStartPacket{})
}

// StopTestSigns broadcasts the test signs stop command
func (c *HanoverController) StopTestSigns() error {
	return c.writeAndRead(packet.TestSignsStopPacket{})
}

func (c *HanoverController) DrawImage(img *image.Gray, signName string) error {
	sign, ok := c.signs[signName]
	if !ok {
		return errors.New("sign not found")
	}

	if err := sign.ValidateImage(img); err != nil {
		return fmt.Errorf("invalid image: %w", err)
	}

	pkt := packet.ImagePacket{
		Address: sign.Address,
		Image:   img,
	}

	bytes, err := pkt.GetBytes()
	if err != nil {
		return fmt.Errorf("failed to get packet bytes: %w", err)
	}

	log.Printf("Sending packet: %X", bytes)
	_, err = c.port.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	return nil
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
	log.Printf("Sending packet: %s", hex.EncodeToString(bytes))
	n, err := c.port.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}
	if n != len(bytes) {
		return fmt.Errorf("incomplete write: wrote %d bytes out of %d", n, len(bytes))
	}
	log.Printf("Wrote %d bytes to serial port", n)
	return nil
}

func (c *HanoverController) writeAndRead(pkt packet.Packet) error {
	if err := c.write(pkt); err != nil {
		return err
	}

	// Read response with timeout
	buf := make([]byte, 128)
	readChan := make(chan readResult)
	go func() {
		n, err := c.port.Read(buf)
		readChan <- readResult{n: n, err: err}
	}()

	select {
	case result := <-readChan:
		if result.err != nil && !errors.Is(result.err, io.EOF) {
			log.Printf("Failed to read response: %v", result.err)
		} else if result.n > 0 {
			log.Printf("Received response: %s", hex.EncodeToString(buf[:result.n]))
		} else {
			log.Println("No data received from read operation")
		}
	case <-time.After(2 * time.Second):
		log.Println("Read operation timed out after 2 seconds")
	}

	return nil
}

type readResult struct {
	n   int
	err error
}
