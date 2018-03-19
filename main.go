package main

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/tonychee7000/Arremi/frontend"
)

func main() {
	defer func() {
		fmt.Println("Arremi Exit.")
	}()
	err := ui.Main(frontend.WindowMain)
	if err != nil {
		fmt.Print("Got error:", err, "\n")
	}
}
