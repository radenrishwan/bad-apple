package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"time"

	badapple "github.com/radenrishwan/bad-apple"
	"golang.org/x/term"
)

const (
	FPS            = 30
	FRAME_DURATION = time.Second / FPS
	FRAME_PATH     = "frames/"
	AUDIO_PATH     = "audio.mp3"
)

func main() {
	// get terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("error getting terminal size:", err)
	}

	// count the frames
	fc, err := os.ReadDir(FRAME_PATH)
	if err != nil {
		log.Fatalln("error while counting frame: ", err)
	}

	frameCount := len(fc)
	for i := 0; i < frameCount; i++ {
		// get the frame path

		framePath := fmt.Sprintf("%s%s", FRAME_PATH, fc[i].Name())

		frame, err := badapple.ProcessFrame(framePath)
		if err != nil {
			log.Fatalln(err)
		}

		// scale the frame
		frame = scaleFrame(frame, width, height)

		renderFrame(frame)

		time.Sleep(FRAME_DURATION)
	}
}

func scaleFrame(frame [][]uint8, targetWidth, targetHeight int) [][]uint8 {
	originalHeight := len(frame)
	originalWidth := len(frame[0])

	scaled := make([][]uint8, targetHeight)
	for i := range scaled {
		scaled[i] = make([]uint8, targetWidth)
	}

	// scale factors
	scaleY := float64(originalHeight) / float64(targetHeight)
	scaleX := float64(originalWidth) / float64(targetWidth)

	// scale the frame
	for y := 0; y < targetHeight; y++ {
		sourceY := int(float64(y) * scaleY)
		if sourceY >= originalHeight {
			sourceY = originalHeight - 1
		}

		for x := 0; x < targetWidth; x++ {
			sourceX := int(float64(x) * scaleX)
			if sourceX >= originalWidth {
				sourceX = originalWidth - 1
			}

			scaled[y][x] = frame[sourceY][sourceX]
		}
	}

	return scaled
}

func renderFrame(frame [][]uint8) {
	os.Stdout.Write([]byte("\033[H"))

	for _, row := range frame {
		for _, pixel := range row {
			if pixel == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("█")
			}
		}

		fmt.Println()
	}
}
