package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
)

func main() {
	// Create a controller
	ctrl, err := goflipdot.NewController("/dev/ttyUSB0")
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	// Add a sign
	if err := ctrl.AddSign("dev", 1, 86, 7, false); err != nil {
		log.Fatalf("Failed to add sign: %v", err)
	}

	// Start the test sequence
	if err := ctrl.StartTestSigns(); err != nil {
		log.Fatalf("Failed to start test signs: %v", err)
	}

	fmt.Println("Test sequence started. Waiting for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Stop the test sequence
	if err := ctrl.StopTestSigns(); err != nil {
		log.Fatalf("Failed to stop test signs: %v", err)
	}

	// Create a 'checkerboard' image
	img, err := ctrl.CreateImage("dev")
	if err != nil {
		log.Fatalf("Failed to create image: %v", err)
	}

	log.Printf("Created image with dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

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
		log.Fatalf("Failed to draw image: %v", err)
	}

	fmt.Println("Checkerboard image drawn. Press Enter to exit...")
	fmt.Scanln()
}
