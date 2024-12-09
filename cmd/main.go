package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"
	"time"

	badapple "github.com/radenrishwan/bad-apple"
)

const (
	FPS            = 30
	FRAME_DURATION = time.Second / FPS
	FRAME_PATH     = "frames/"
	AUDIO_PATH     = "audio.mp3"
)

func main() {
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

		// render frame
		renderFrame(frame)

		time.Sleep(FRAME_DURATION)
	}
}

func renderFrame(frame [][]uint8) {
	fmt.Print("\033[H\033[2J")

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
