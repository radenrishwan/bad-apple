package main

import (
	_ "image/png"

	badapple "github.com/radenrishwan/bad-apple"
)

const (
	FPS        = 30
	FRAME_PATH = "frames/"
	AUDIO_PATH = "bad-apple.mp3"
)

func main() {
	badapple.RunCLI(FPS, FRAME_PATH, AUDIO_PATH)
}
