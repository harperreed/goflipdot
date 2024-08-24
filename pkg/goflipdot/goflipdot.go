package goflipdot

import (
	"github.com/harperreed/goflipdot/internal/controller"
	"github.com/harperreed/goflipdot/internal/sign"
	"image"
	"io"
)

// Controller represents the main interface for controlling Hanover flipdot displays
type Controller struct {
	ctrl *controller.HanoverController
}

// NewController creates a new Controller
func NewController(port io.Writer) *Controller {
	return &Controller{
		ctrl: controller.NewHanoverController(port),
	}
}

// AddSign adds a new sign to the controller
func (c *Controller) AddSign(name string, address, width, height int, flip bool) error {
	s := sign.NewHanoverSign(address, width, height, flip)
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
		return nil, err
	}
	return s.CreateImage(), nil
}
