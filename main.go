package main

import (
	"fmt"

	"github.com/Darrell-Devana/be-note-app/util"
)

func main() {
	message := "Hello this is a message"
	messageUpper := util.ToUpperCase(message)
	fmt.Println(messageUpper)
}
