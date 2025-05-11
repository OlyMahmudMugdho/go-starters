package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd1 := exec.Command("curl", "https://jsonplaceholder.typicode.com/todos/11")
	cmd2 := exec.Command("grep", "title")

	// Create the pipe
	pipe, err := cmd1.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd2.Stdin = pipe

	// Start the second command first
	cmd2.Stdout = os.Stdout // Or capture output differently
	if err := cmd2.Start(); err != nil {
		log.Fatal(err)
	}

	// Then run the first command
	if err := cmd1.Run(); err != nil {
		log.Fatal(err)
	}

	// Wait for the second command to complete
	if err := cmd2.Wait(); err != nil {
		// grep returns error when no matches found, which is normal
		if exiterr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("grep exited with status %d\n", exiterr.ExitCode())
		} else {
			log.Fatal(err)
		}
	}
}
