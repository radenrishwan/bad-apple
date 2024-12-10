package badapple

import (
	"errors"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func ProcessFrame(path string) ([][]uint8, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Error while opening file :" + err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.New("Error while decoding image :" + err.Error())
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	binaryArray := make([][]uint8, height)
	for y := 0; y < height; y++ {
		binaryArray[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// get the pixel color
			pixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray)

			// applying threshold
			var value uint8
			if pixel.Y < 128 {
				value = 0
			} else {
				value = 1
			}
			binaryArray[y][x] = value
		}
	}

	return binaryArray, nil
}

func LoadFrame(path string, targetHeight, targetWidth int) ([][]uint8, error) {
	frame, err := ProcessFrame(path)
	if err != nil {
		return nil, err
	}

	// get original size
	originalHeight := len(frame)
	originalWidth := len(frame[0])

	// create a new array for the resized frame
	resized := make([][]uint8, targetHeight)
	for i := range resized {
		resized[i] = make([]uint8, targetWidth)
	}

	// resize using nearest neighbor scaling
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// calculate corresponding position in original image
			sourceY := y * originalHeight / targetHeight
			sourceX := x * originalWidth / targetWidth

			resized[y][x] = frame[sourceY][sourceX]
		}
	}

	return resized, nil
}

func PlayAudio(path string, done <-chan bool) error {
	audio, err := os.Open(path)
	defer audio.Close()

	streamer, format, err := mp3.Decode(audio)
	if err != nil {
		return errors.New("Error while decoding audio :" + err.Error())
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return errors.New("Error while initializing speaker :" + err.Error())
	}

	speaker.Play(streamer)

	<-done

	return nil
}
