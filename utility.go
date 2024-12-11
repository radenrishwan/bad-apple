package badapple

import (
	"errors"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var VOLUME_LEVEL = -2.0

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

	grayArray := make([][]uint8, height)
	for y := 0; y < height; y++ {
		grayArray[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// get the pixel color
			pixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray)

			// convert to 10 levels (0-9)
			switch {
			case pixel.Y < 51:
				grayArray[y][x] = 0
			case pixel.Y < 102:
				grayArray[y][x] = 1
			case pixel.Y < 153:
				grayArray[y][x] = 2
			case pixel.Y < 204:
				grayArray[y][x] = 3
			default:
				grayArray[y][x] = 4
			}
		}
	}

	return grayArray, nil
}

func LoadFrame(path string, targetHeight, targetWidth int) ([][]uint8, error) {
	frame, err := ProcessFrame(path)
	if err != nil {
		return nil, err
	}

	if targetWidth == 0 && targetHeight == 0 {
		return frame, nil
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

	// Create volume controller with the flag value
	volumeStreamer := &effects.Volume{
		Streamer: streamer,
		Volume:   VOLUME_LEVEL,
		Base:     2,
	}

	speaker.Play(volumeStreamer)

	<-done

	return nil
}
