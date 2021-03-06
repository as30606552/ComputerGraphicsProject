package pngimage

import (
	"image/color"
	"math/rand"
)

// A structure for storing colors in RGB format without specifying alfa value.
// All pixels have a maximum alfa value, meaning they are completely opaque.
// Implements the interface color.Color, so that all the functions that work with color can be used.
type RGB struct {
	R, G, B uint8
}

// Implementation of the RGBA method in the color.Color interface.
func (rgb RGB) RGBA() (r, g, b, a uint32) {
	return uint32(rgb.R), uint32(rgb.G), uint32(rgb.B), 255
}

// Converts an RGB object to an color.RGBA object.
func (rgb RGB) ToRGBA() color.RGBA {
	return color.RGBA{
		R: rgb.R,
		G: rgb.G,
		B: rgb.B,
		A: 255,
	}
}

// Creates black RGB color.
func BlackColor() RGB {
	return RGB{R: 0, G: 0, B: 0}
}

// Creates white RGB color.
func WhiteColor() RGB {
	return RGB{R: 255, G: 255, B: 255}
}

// Creates red RGB color.
func RedColor() RGB {
	return RGB{R: 255, G: 0, B: 0}
}

// Creates green RGB color.
func GreenColor() RGB {
	return RGB{R: 0, G: 255, B: 0}
}

// Creates blue RGB color.
func BlueColor() RGB {
	return RGB{R: 0, G: 0, B: 255}
}

// Creates random RGB color.
func RandomColor() RGB {
	return RGB{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
	}
}
