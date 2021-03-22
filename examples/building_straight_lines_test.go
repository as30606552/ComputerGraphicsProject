package examples

import (
	"computer_graphics/pngimage"
	"fmt"
	"image"
	"math"
)

// The type for passing a method to a function starLine.
type drawingStraightLine func(point1, point2 image.Point, img *pngimage.Image, rgb *pngimage.RGB)

// Creates a 12-pointed star line using the specified method.
func starLine(method drawingStraightLine, img *pngimage.Image, rgb *pngimage.RGB) {
	for i := 0; i < 12; i++ {
		alpha := (float64(2*i) * math.Pi) / 13
		x := int(100 + 95*math.Cos(alpha))
		y := int(100 + 95*math.Sin(alpha))
		method(image.Point{100, 100}, image.Point{x, y}, img, rgb)
	}
}

// Create a line by drawing N points on a straight line.
func SimplestMethod(point1, point2 image.Point, img *pngimage.Image, rgb *pngimage.RGB) {
	for t := 0.0; t < 1.0; t += 0.01 {
		x := int(float64(point1.X)*(1.0-t) + float64(point2.X)*t)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		img.Set(x, y, *rgb)
	}
}

// Create a line by set the x values in increments of one pixel, and calculate the y values for them.
func SecondMethod(point1, point2 image.Point, img *pngimage.Image, rgb *pngimage.RGB) {
	for x := point1.X; x <= point2.X; x++ {
		t := float64(x-point1.X) / float64(point2.X-point1.X)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		img.Set(x, y, *rgb)
	}
}

// Create a line by an improved the third method.
func ThirdMethod(point1, point2 image.Point, img *pngimage.Image, rgb *pngimage.RGB) {
	steep := false
	if math.Abs(float64(point1.X-point2.X)) < math.Abs(float64(point1.Y-point2.Y)) {
		point1.X, point1.Y = point1.Y, point1.X
		point2.X, point2.Y = point2.Y, point2.X
		steep = true
	}
	if point1.X > point2.X {
		point1.X, point2.X = point2.X, point1.X
		point1.Y, point2.Y = point2.Y, point1.Y
	}
	for x := point1.X; x <= point2.X; x++ {
		t := float64(x-point1.X) / float64(point2.X-point1.X)
		y := int(float64(point1.Y)*(1.0-t) + float64(point2.Y)*t)
		if steep {
			img.Set(y, x, *rgb)
		} else {
			img.Set(x, y, *rgb)
		}
	}
}

// Example of creating a simplest method image.
func ExampleSimplestMethod() {
	var img = pngimage.WhiteImage(200, 200)
	var rgb = pngimage.NewRGB(255, 0, 0)
	starLine(SimplestMethod, &img, rgb)
	img.Save("pictures/simplest_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a second method image.
func ExampleSecondMethod() {
	var img = pngimage.WhiteImage(200, 200)
	var rgb = pngimage.NewRGB(255, 0, 0)
	starLine(SecondMethod, &img, rgb)
	img.Save("pictures/second_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}

// Example of creating a third method image.
func ExampleThirdMethod() {
	var img = pngimage.WhiteImage(200, 200)
	var rgb = pngimage.NewRGB(255, 0, 0)
	starLine(ThirdMethod, &img, rgb)
	img.Save("pictures/third_method_image.png")
	fmt.Println("Ok")
	// Output: Ok
}
