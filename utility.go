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
