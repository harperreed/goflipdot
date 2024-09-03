package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/harperreed/goflipdot/pkg/goflipdot"
)

const (
	signAddress = 1
	signColumns = 96
	signRows    = 16
)

func main() {
	serialPort := flag.String("port", "/dev/pts/7", "Serial port for the flipdot display")
	patternNum := flag.Int("pattern", -1, "Pattern number to display (0-5), or -1 for all patterns")
	flag.Parse()

	if *serialPort == "" {
		log.Fatal("Serial port must be specified")
	}

	ctrl, err := goflipdot.NewController(*serialPort)
	if err != nil {
		log.Fatal(err)
	}

	if err := ctrl.AddSign("dev", signAddress, signColumns, signRows); err != nil {
		log.Fatal(err)
	}

	patterns := GetPatterns()
	patternNames := []string{
		"1s at row edges",
		"1s on borders",
		"Checkerboard",
		"All pixels on",
		"Alternating columns",
		"Large 'X'",
		"Clear",
	}

	if *patternNum >= 0 && *patternNum < len(patternNames) {
		// Display only the specified pattern
		name := patternNames[*patternNum]
		patternFunc := patterns[name]
		displayPattern(ctrl, name, patternFunc)
	} else if *patternNum == -1 {
		// Display all patterns
		for _, name := range patternNames {
			patternFunc := patterns[name]
			displayPattern(ctrl, name, patternFunc)
			time.Sleep(2 * time.Second)
		}
	} else {
		fmt.Println("Invalid pattern number. Use -1 for all patterns or 0-5 for a specific pattern.")
		os.Exit(1)
	}

	fmt.Println("Test sequence completed.")
}

func displayPattern(ctrl *goflipdot.Controller, name string, patternFunc Pattern) {
	img := patternFunc(signColumns, signRows)
	printArrayInfo(img, name)
	err := ctrl.DrawImage(img, "dev")
	if err != nil {
		log.Printf("Failed to draw image: %v", err)
	}
}

func printArrayInfo(img *image.Gray, name string) {
	fmt.Printf("\n%s:\n", name)
	fmt.Printf("Shape: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
	fmt.Print("First row: ")
	for x := 0; x < img.Bounds().Dx(); x++ {
		if img.GrayAt(x, 0).Y > 0 {
			fmt.Print("1 ")
		} else {
			fmt.Print("0 ")
		}
	}
	fmt.Println()
	fmt.Print("Last row: ")
	for x := 0; x < img.Bounds().Dx(); x++ {
		if img.GrayAt(x, img.Bounds().Dy()-1).Y > 0 {
			fmt.Print("1 ")
		} else {
			fmt.Print("0 ")
		}
	}
	fmt.Println()
	sum := 0
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.GrayAt(x, y).Y > 0 {
				sum++
			}
		}
	}
	fmt.Printf("Sum of elements: %d\n", sum)
}
