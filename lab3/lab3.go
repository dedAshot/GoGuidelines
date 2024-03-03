package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

const scale = 3 // default scale = 1 / 2.2

func correctGamma(color uint32) uint8 {
	res := math.Sqrt(float64(color) * scale)

	if res < 0 {
		res = 0
	} else if res > 255 {
		res = 255
	}
	return uint8(res)
}

func handleError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

// function wrapper, print execution time
func bench(f func(), msg string) {
	timeStart := time.Now()
	f()
	timeEnd := time.Now()
	fmt.Println(msg, "in", timeEnd.Sub(timeStart))
}

func ReadFile(byteArr *[]byte, path string) {
	file, err := os.Open(path)
	handleError("file open error:", err)
	defer file.Close()
	fmt.Println("File opened")

	*byteArr, err = io.ReadAll(file)
	handleError("file read error:", err)
	fmt.Println("File readed")
}

func DecodeFile(imageRef *draw.Image, fileBuff *[]byte) {
	img, _, err := image.Decode(bytes.NewReader(*fileBuff))
	handleError("Decode error: ", err)
	fmt.Println("File decoded")

	var ok bool
	*imageRef, ok = img.(draw.Image)
	if !ok {
		handleError("Cant cast image.Image to draw.Image ", nil)
	}
	img = nil
}

func chunkGammaCorrection(image *draw.Image, Xmin, Ymin, Xmax, Ymax int) {

	for y := Ymin; y < Ymax; y++ {
		for x := Xmin; x < Xmax; x++ {
			r, g, b, a := (*image).At(x, y).RGBA()
			(*image).Set(x, y, color.RGBA{correctGamma(r), correctGamma(g), correctGamma(b), uint8(a)})
		}
	}
}

func correctImageGamma(image *draw.Image, chunkSize uint64) {
	CORES := runtime.NumCPU() //logical cores (cpu threads)

	waitGroup := &sync.WaitGroup{}

	waitGroupInputChan := make(chan struct {
		image                  *draw.Image
		Xmin, Ymin, Xmax, Ymax int
	}, CORES)

	Ymin := (*image).Bounds().Min.Y
	Ymax := (*image).Bounds().Max.Y

	Xmin := (*image).Bounds().Min.X
	Xmax := (*image).Bounds().Max.X

	//send tasks
	go func() {
		for y := Ymin; y < Ymax; y += int(chunkSize) {
			for x := Xmin; x < Xmax; x += int(chunkSize) {
				xChunkMax := int(math.Min(float64(x+int(chunkSize)), float64(Xmax)))
				yChunkMax := int(math.Min(float64(y+int(chunkSize)), float64(Ymax)))

				waitGroupInputChan <- struct {
					image                  *draw.Image
					Xmin, Ymin, Xmax, Ymax int
				}{
					image,
					x, y, xChunkMax, yChunkMax,
				}
			}
		}
		close(waitGroupInputChan)
	}()

	//create workers
	for i := 0; i < CORES; i++ {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for cort := range waitGroupInputChan {
				chunkGammaCorrection(cort.image, cort.Xmin, cort.Ymin, cort.Xmax, cort.Ymax)
			}
		}()
	}

	waitGroup.Wait()
}

func EncodeFile(fileName string, image *draw.Image) {
	file, err := os.Create(fileName)
	handleError("Unable to create file:", err)
	defer file.Close()

	err = png.Encode(file, *image)
	handleError("Encode error: ", err)
}

func main() {
	fmt.Println("Hello, World!")
	defer fmt.Println("Goodbye World")

	var fileBuff []byte
	var resImage draw.Image

	bench(func() {
		ReadFile(&fileBuff, "./sourceImages/heic1502a.png")
	}, "ReadFile")

	bench(func() {
		DecodeFile(&resImage, &fileBuff)
	}, "Decode file")

	//clear memory
	fileBuff = nil

	bench(func() {
		correctImageGamma(&resImage, 4096)
	}, "Correct gamma")

	bench(func() {
		EncodeFile("./outputImages/resImg.png", &resImage)
	}, "File Encode")
}
