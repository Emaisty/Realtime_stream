package main

import (
	"fmt"
	"os"

	"screenshot"
)

func main() {
	screenshot.InitConn()
	id := screenshot.GetActiveWindow()
	fmt.Printf("%v", id)
	screenshot.CloseConn()
	os.Exit(0)
}
