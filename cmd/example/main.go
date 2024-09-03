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
	serialPort := "/dev/pts/7" // Update this to match your system
	log.Printf("Initializing Hanover Controller with serial port: %s", serialPort)
	ctrl, err := controller.NewHanoverController(serialPort)
	if err != nil {
		log.Fatalf("Failed to create Hanover Controller: %v", err)
	}
	log.Println("Hanover Controller created successfully")

	log.Println("Creating new Hanover Sign")
	sign, err := sign.NewHanoverSign(1, 96, 16) // Update dimensions to match your sign
	if err != nil {
		log.Fatalf("Failed to create Hanover Sign: %v", err)
	}
	log.Printf("Hanover Sign created with dimensions: 96x16")

	log.Println("Adding sign to controller")
	if err := ctrl.AddSign("dev", sign); err != nil {
		log.Fatalf("Failed to add sign to controller: %v", err)
	}
	log.Println("Sign added successfully")

	// Create a test image
	log.Println("Creating test image")
	img := sign.CreateImage()
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.White)
				log.Printf("Set pixel at (%d, %d) to White", x, y)
			} else {
				img.Set(x, y, color.Black)
				log.Printf("Set pixel at (%d, %d) to Black", x, y)
			}
		}
	}
	log.Println("Test image created successfully")

	// Draw the image
	log.Println("Sending image to flipdot display")
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Fatalf("Failed to draw image: %v", err)
	}
	log.Println("Image sent successfully")

	fmt.Println("Image sent to flipdot display. Press Enter to exit...")
	log.Println("Waiting for 5 seconds before exit")
	time.Sleep(5 * time.Second)
	log.Println("Exiting program")
}
