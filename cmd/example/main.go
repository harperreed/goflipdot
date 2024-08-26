package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
	"github.com/tarm/serial"
)

func main() {
	config := &serial.Config{
		Name:     "/dev/ttyUSB0",
		Baud:     4800,
		Size:     8,
		Parity:   serial.ParityNone,
		StopBits: serial.Stop1,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	ctrl, err := goflipdot.NewController(port)
	if err != nil {
		log.Fatal(err)
	}

	if err := ctrl.AddSign("dev", 1, 86, 7, false); err != nil {
		log.Fatal(err)
	}

	// Test sequence
	log.Println("Starting test sequence")
	if err := ctrl.StartTestSigns(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(10 * time.Second)

	log.Println("Stopping test sequence")
	if err := ctrl.StopTestSigns(); err != nil {
		log.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	// All pixels on
	log.Println("Setting all pixels on")
	img, _ := ctrl.CreateImage("dev")
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.White)
		}
	}
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Printf("Error drawing all-on image: %v", err)
	}
	time.Sleep(5 * time.Second)

	// All pixels off
	log.Println("Setting all pixels off")
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, color.Black)
		}
	}
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Printf("Error drawing all-off image: %v", err)
	}
	time.Sleep(5 * time.Second)

	// Single column on
	log.Println("Setting single column on")
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if x == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
	if err := ctrl.DrawImage(img, "dev"); err != nil {
		log.Printf("Error drawing single column: %v", err)
	}
	time.Sleep(5 * time.Second)

	log.Println("Program completed. Press Enter to exit...")
	fmt.Scanln()
}
