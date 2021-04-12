package main

import (
	"fmt"

	"github.com/BrenekH/encodarr/runner/cmd_runner"
)

func main() {
	cR := cmd_runner.NewCmdRunner()

	cR.Start([]string{"-i", "/home/brenekh/Downloads/2Fast2Furious.mp4", "/home/brenekh/out.mp4"})

	for {
		if cR.Done() {
			fmt.Println("Command Runner done")
			break
		}
	}

	fmt.Println("Program exiting")
}
