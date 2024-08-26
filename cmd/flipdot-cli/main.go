package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	"github.com/harperreed/goflipdot/internal"
)

func main() {
	portName := flag.String("port", "/dev/ttyUSB0", "Serial port name")
	command := flag.String("cmd", "", "Command to send (start_test, stop_test, draw_pattern, or send_byte)")
	byteToSend := flag.Int("byte", 0xFF, "Byte to send when using send_byte command")
	verbose := flag.Bool("v", false, "Verbose mode")
	flag.Parse()

	config := &serial.Config{
		Name:        *portName,
		Baud:        internal.BaudRate,
		ReadTimeout: time.Second * 5,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	if *verbose {
		fmt.Printf("Opening port %s with config: %+v\n", *portName, config)
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer port.Close()

	var packet []byte

	switch *command {
	case "start_test":
		packet = internal.FormatPacket(internal.CommandCodes["start_test_signs"], '0', nil)
	case "stop_test":
		packet = internal.FormatPacket(internal.CommandCodes["stop_test_signs"], '0', nil)
	case "draw_pattern":
		imageData := []byte{0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA}
		payload := append([]byte{byte(len(imageData)), 0}, internal.ToAsciiHex(imageData)...)
		packet = internal.FormatPacket(internal.CommandCodes["write_image"], '1', payload)
	case "send_byte":
		packet = []byte{byte(*byteToSend)}
	default:
		log.Fatalf("Unknown command: %s", *command)
	}

	if *verbose {
		fmt.Printf("Sending packet: %X\n", packet)
		fmt.Printf("ASCII representation: %s\n", string(packet))
	}
	n, err := port.Write(packet)
	if err != nil {
		log.Fatalf("Failed to write to serial port: %v", err)
	}
	if *verbose {
		fmt.Printf("Sent %d bytes\n", n)
	}

	// Add a delay to match the 2 Hz refresh rate
	time.Sleep(500 * time.Millisecond)

	// Try to read any response
	buf := make([]byte, 128)
	n, err = port.Read(buf)
	if err != nil {
		if *verbose {
			fmt.Printf("No response received (this may be normal): %v\n", err)
		}
	} else {
		fmt.Printf("Received response: %X\n", buf[:n])
		fmt.Printf("ASCII representation: %s\n", string(buf[:n]))
	}

	fmt.Println("Command completed")
}
