package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/diegopacheco/writing-interpreter-in-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
