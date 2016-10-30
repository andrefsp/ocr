package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

var ALPHABET = map[string]string{
	"A": "./alphabet/A_50x75.png",
	"B": "./alphabet/B_50x75.png",
}

func SaveImage(img image.Image, filePath string) error {
	newFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	return png.Encode(newFile, img)
}

func OpenImageFile(filename string) image.Image {

	fileReader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(fileReader)
	if err != nil {
		panic(err)
	}

	return img
}

func GetImageBounds(img image.Image) (image.Point, image.Point) {
	var minX, maxX, minY, maxY int

	perdominantColor := img.At(0, 0)

	minPoint := img.Bounds().Min
	maxPoint := img.Bounds().Max

	for x := minPoint.X; x < maxPoint.X; x++ {
		for y := minPoint.Y; y < maxPoint.Y; y++ {
			if img.At(x, y) != perdominantColor {
				if minX == 0 || x < minX {
					minX = x
				}
				if maxX == 0 || x > maxX {
					maxX = x
				}
				if minY == 0 || y < minY {
					minY = y
				}
				if maxY == 0 || y > maxY {
					maxY = y
				}
			}
		}
	}

	return image.Point{X: minX, Y: minY}, image.Point{X: maxX, Y: maxY}
}

// Very simply Eucledean distance where it encodes black and white as 0 or 1
func EuclideanDistance(baseImage image.Image, compareImage image.Image) float64 {

	maxPoint := baseImage.Bounds().Max

	var sum float64

	p := float64(2)

	for x := 1; x < maxPoint.X; x++ {
		for y := 1; y < maxPoint.Y; y++ {

			baseR, baseG, baseB, baseA := baseImage.At(x, y).RGBA()
			compareR, compareG, compareB, compareA := compareImage.At(x, y).RGBA()

			sum += math.Pow(float64(baseR-compareR), p) + math.Pow(float64(baseG-compareG), p) +
				math.Pow(float64(baseB-compareB), p) + math.Pow(float64(baseA-compareA), p)

		}
	}
	return math.Sqrt(sum)
}

func CropImage(baseImage image.Image, minPoint image.Point, maxPoint image.Point) image.Image {
	targetImage := image.NewNRGBA(image.Rectangle{
		Min: image.Point{
			X: 1,
			Y: 1,
		},
		Max: image.Point{
			X: maxPoint.X - minPoint.X,
			Y: maxPoint.Y - minPoint.Y,
		},
	})

	for x := minPoint.X; x < maxPoint.X; x++ {
		for y := minPoint.Y; y <= maxPoint.Y; y++ {
			color := baseImage.At(x, y)
			targetImage.Set(x-minPoint.X, y-minPoint.Y, color)
		}
	}

	return targetImage
}

// ResizeImage simply does a bilinear intrepolation on from the main image
// and returns the new scaled image.
func ResizeImage(baseImage image.Image, width int, height int) image.Image {

	newImage := image.NewNRGBA(image.Rectangle{
		Min: image.Point{
			X: 1,
			Y: 1,
		},
		Max: image.Point{
			X: width,
			Y: height,
		},
	})
	maxPoint := baseImage.Bounds().Max

	XAxisRatio := float32(maxPoint.X) / float32(width)
	YAxixRatio := float32(maxPoint.Y) / float32(height)

	for x := 1; x < width; x++ {
		for y := 1; y < height; y++ {

			fromX := int(XAxisRatio * float32(x))
			fromY := int(YAxixRatio * float32(y))

			newImage.Set(x, y, baseImage.At(fromX, fromY))
		}
	}

	return newImage
}

func main() {

	img := OpenImageFile("./images/B.png")

	alphabetImage := OpenImageFile("./alphabet/B_50x75.png")

	minPoint, maxPoint := GetImageBounds(img)

	croppedImage := CropImage(img, minPoint, maxPoint)

	resizedImage := ResizeImage(
		croppedImage,
		alphabetImage.Bounds().Max.X,
		alphabetImage.Bounds().Max.Y)

	for letter, filename := range ALPHABET {
		compareImage := OpenImageFile(filename)
		fmt.Printf("%s :  %f \n", letter, EuclideanDistance(resizedImage, compareImage))
	}

}
