package goflipdot

import (
	"fmt"
	"image"

	"github.com/harperreed/goflipdot/internal/controller"
	"github.com/harperreed/goflipdot/internal/sign"
)

// Controller represents the main interface for controlling Hanover flipdot displays
type Controller struct {
	ctrl *controller.HanoverController
}

// NewController creates a new Controller
func NewController(serialPort string) (*Controller, error) {
	ctrl, err := controller.NewHanoverController(serialPort)
	if err != nil {
		return nil, fmt.Errorf("failed to create controller: %w", err)
	}
	return &Controller{
		ctrl: ctrl,
	}, nil
}

// AddSign adds a new sign to the controller
func (c *Controller) AddSign(name string, address, width, height int) error {
	s, err := sign.NewHanoverSign(address, width, height)
	if err != nil {
		return fmt.Errorf("failed to create sign: %w", err)
	}
	return c.ctrl.AddSign(name, s)
}

// StartTestSigns starts the test sequence on all connected signs
func (c *Controller) StartTestSigns() error {
	return c.ctrl.StartTestSigns()
}

// StopTestSigns stops the test sequence on all connected signs
func (c *Controller) StopTestSigns() error {
	return c.ctrl.StopTestSigns()
}

// DrawImage sends an image to a specific sign
func (c *Controller) DrawImage(img *image.Gray, signName string) error {
	return c.ctrl.DrawImage(img, signName)
}

// CreateImage creates a blank image for a specific sign
func (c *Controller) CreateImage(signName string) (*image.Gray, error) {
	s, err := c.ctrl.GetSign(signName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sign: %w", err)
	}
	return s.CreateImage(), nil
}
