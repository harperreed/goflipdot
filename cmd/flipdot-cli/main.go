package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
	"github.com/tarm/serial"
)

func main() {
	portName := flag.String("port", "/dev/ttyUSB0", "Serial port name")
	command := flag.String("cmd", "", "Command to send (start_test, stop_test, draw_pattern, or send_byte)")
	byteToSend := flag.Int("byte", 0xFF, "Byte to send when using send_byte command")
	verbose := flag.Bool("v", false, "Verbose mode")
	flag.Parse()

	config := &serial.Config{
		Name:        *portName,
		Baud:        goflipdot.BaudRate,
		ReadTimeout: time.Second * 5,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	if *verbose {
		fmt.Printf("Opening port %s with config: %+v\n", *portName, config)
	}

	controller, err := goflipdot.NewController(*portName)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	switch *command {
	case "start_test":
		err = controller.StartTestSigns()
	case "stop_test":
		err = controller.StopTestSigns()
	case "draw_pattern":
		// Create a simple pattern (you may need to adjust this based on your sign dimensions)
		img := goflipdot.CreateCheckerboardPattern(86, 7)
		err = controller.DrawImage(img, "")
	case "send_byte":
		// This functionality might not be directly available in the goflipdot package
		// You may need to implement a custom method or use a different approach
		log.Println("send_byte command is not supported in this version")
	default:
		log.Fatalf("Unknown command: %s", *command)
	}

	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	if *verbose {
		fmt.Println("Command completed successfully")
	}
}
