package main

import (
	_ "image/png"

	badapple "github.com/radenrishwan/bad-apple"
)

const (
	FPS        = 30
	FRAME_PATH = "frames/"
	AUDIO_PATH = "bad-apple.mp3"
	HEIGHT     = 32 // DEFAULT is 0, its mean no scale
	WIDTH      = 32
)

func main() {
	badapple.RunCLI(FPS, FRAME_PATH, AUDIO_PATH, HEIGHT, WIDTH)
}
