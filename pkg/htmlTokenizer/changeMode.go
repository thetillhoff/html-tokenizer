package tokenizer

import "log"

func changeMode(targetMode Mode) Mode {
	log.Println("Switching from " + mode.String() + " to " + targetMode.String() + " mode")
	word = "" // Reset word
	return targetMode
}
