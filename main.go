package main

import (
	"fmt"
	"github.com/opticaline/mark-favorite-command/chooser"
	"github.com/opticaline/mark-favorite-command/history"
	"log"
	"strings"
)

func main() {
	his, err := history.OpenHistory()
	if err != nil {
		log.Fatalln("Can't get history")
	}
	lines := his.GetHistory()
	commands, err := chooser.Construct(lines).WaitForAnswer()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", strings.Join(commands, "\n"))
}
