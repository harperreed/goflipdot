package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
)

func main() {
	// Open a serial port (this is just a placeholder, you'd need to use a real serial library)
	port, err := os.OpenFile("/dev/ttyUSB0", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// Create a controller
	ctrl, err := goflipdot.NewController(port)
	if err != nil {
		log.Fatal(err)
	}

	// Add a sign
	if err := ctrl.AddSign("dev", 1, 86, 7, false); err != nil {
		log.Fatal(err)
	}

	// Start the test sequence
	if err := ctrl.StartTestSigns(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Test sequence started. Press Enter to stop...")
	fmt.Scanln()

	// Stop the test sequence
	if err := ctrl.StopTestSigns(); err != nil {
		log.Fatal(err)
	}

	// Create a 'checkerboard' image
	img, err := ctrl.CreateImage("dev")
	if err != nil {
		log.Fatal(err)
	}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.White)
			}
		}
	}

	// Draw the image
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Checkerboard image drawn. Press Enter to exit...")
	fmt.Scanln()
}
