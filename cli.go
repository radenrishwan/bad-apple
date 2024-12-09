package badapple

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

func RunCLI(fps int, framePath, audioPath string) {
	FRAME_DURATION := int(time.Second) / fps

	// get terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Fatalln("error getting terminal size:", err)
	}

	// count the frames
	fc, err := os.ReadDir(framePath)
	if err != nil {
		log.Fatalln("error while counting frame: ", err)
	}

	// start the audio
	done := make(chan bool)
	go PlayAudio(audioPath, done)

	// sorting the frames
	sort.Slice(fc, func(i, j int) bool {
		// sort by frame number
		// e.g. frame_1.png, frame_2.png, frame_3.png, ...
		a := strings.Split(fc[i].Name(), "_")
		b := strings.Split(fc[j].Name(), "_")

		// remove the .png
		a = strings.Split(a[1], ".")
		b = strings.Split(b[1], ".")

		numA, _ := strconv.Atoi(a[0])
		numB, _ := strconv.Atoi(b[0])

		return numA < numB
	})

	// start the frame
	frameCount := len(fc)
	for i := 0; i < frameCount; i++ {
		cycleStart := time.Now()

		// get the frame path
		framePath := fmt.Sprintf("%s%s", framePath, fc[i].Name())

		frame, err := ProcessFrame(framePath)
		if err != nil {
			log.Fatalln(err)
		}

		// scale the frame
		frame = scaleFrame(frame, width, height)

		// render the frame
		renderFrame(frame)
		// fmt.Println("rendering ", framePath)

		processTime := time.Since(cycleStart)
		if processTime < time.Duration(FRAME_DURATION) {
			time.Sleep(time.Duration(FRAME_DURATION) - processTime)
		}
	}

	<-done
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
	// clear the screen and move the cursor to the top left
	os.Stdout.Write([]byte("\033[H\033[2J"))

	for _, row := range frame {
		for _, pixel := range row {
			if pixel == 0 {
				os.Stdout.Write([]byte(" "))
			} else {
				os.Stdout.Write([]byte("█"))
			}
		}

		os.Stdout.Write([]byte("\n"))
	}
}
