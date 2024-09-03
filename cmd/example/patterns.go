package main

import (
	"image"
	"image/color"
)

// Pattern represents a function that generates an image pattern
type Pattern func(width, height int) *image.Gray

// GetPatterns returns a map of pattern names to their generating functions
func GetPatterns() map[string]Pattern {
	return map[string]Pattern{
		"1s at row edges":     createArray1,
		"1s on borders":       createArray2,
		"Checkerboard":        createArray3,
		"All pixels on":       createArray4,
		"Alternating columns": createArray5,
		"Large 'X'":           createArray6,
	}
}

func createArray1(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		img.Set(0, y, color.White)
		img.Set(width-1, y, color.White)
	}
	return img
}

func createArray2(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		img.Set(x, 0, color.White)
		img.Set(x, height-1, color.White)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, color.White)
		img.Set(width-1, y, color.White)
	}
	return img
}

func createArray3(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x+y)%2 == 0 {
				img.Set(x, y, color.White)
			}
		}
	}
	return img
}

func createArray4(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}
	return img
}

func createArray5(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 1; x < width; x += 2 {
			img.Set(x, y, color.White)
		}
	}
	return img
}

func createArray6(width, height int) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for i := 0; i < height; i++ {
		img.Set(i*6, i, color.White)
		img.Set(width-1-i*6, i, color.White)
	}
	return img
}
