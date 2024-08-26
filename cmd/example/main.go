package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/harperreed/goflipdot/internal/controller"
	"github.com/harperreed/goflipdot/internal/sign"
)

func main() {
	serialPort := "/dev/ttyUSB0" // Update this to match your system
	ctrl, err := controller.NewHanoverController(serialPort)
	if err != nil {
		log.Fatal(err)
	}

	sign, err := sign.NewHanoverSign(1, 96, 16) // Update dimensions to match your sign
	if err != nil {
		log.Fatal(err)
	}

	if err := ctrl.AddSign("dev", sign); err != nil {
		log.Fatal(err)
	}

	// Create a test image
	img := sign.CreateImage()
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	// Draw the image
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image sent to flipdot display. Press Enter to exit...")
	time.Sleep(5 * time.Second)
}
