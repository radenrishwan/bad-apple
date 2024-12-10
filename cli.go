package badapple

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func RunCLI(fps int, framePath, audioPath string) {
	FRAME_DURATION := int(time.Second) / fps

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	if err := screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	events := make(chan tcell.Event)
	go func() {
		for {
			ev := screen.PollEvent()
			events <- ev
		}
	}()

	go func() {
		for {
			ev := <-events

			// close the program if the user press CTRL + C or ESC
			switch evt := ev.(type) {
			case *tcell.EventKey:
				if evt.Key() == tcell.KeyCtrlC || evt.Key() == tcell.KeyEscape {
					screen.Fini()
					os.Exit(0)
				}
			}
		}

	}()

	// get terminal size
	// width, height, err := term.GetSize(int(os.Stdout.Fd()))
	// if err != nil {
	// 	log.Fatalln("error getting terminal size:", err)
	// }

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

		frame, err := LoadFrame(framePath+fc[i].Name(), 32, 32)
		if err != nil {
			log.Fatal(err)
		}

		width, height := screen.Size()
		frame = scaleFrame(frame, width, height)

		screen.Clear()
		renderFrameTcell(frame, screen)
		screen.Show()

		processTime := time.Since(cycleStart)
		if processTime < time.Duration(FRAME_DURATION) {
			time.Sleep(time.Duration(FRAME_DURATION) - processTime)
		}
	}

	<-done
}

func renderFrameTcell(frame [][]uint8, screen tcell.Screen) {
	// chars := []rune{' ', '░', '▒', '▓', '█'}

	// alternate characters
	chars := []rune{' ', '.', 'o', 'O', '0'}

	for y, row := range frame {
		for x, pixel := range row {
			style := tcell.StyleDefault
			char := chars[pixel]

			screen.SetContent(x, y, char, nil, style)
		}
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
