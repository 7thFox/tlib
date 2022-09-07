package main

import (
	"fmt"
	"os"
	"strings"
)

const Version = "0.1-alpha"

func main() {
	cmd := "help"
	if len(os.Args) > 1 {
		cmd = strings.ToLower(os.Args[1])
	}

	switch cmd {
	default:
		fmt.Fprintf(os.Stderr, "Unknown command '%s'\n", cmd)
		fallthrough
	case "help":
		fmt.Println("Usage: tlib [command]")
		fmt.Println("Commands:")
		fmt.Println("\t help    \t Print this message")
		fmt.Println("\t version \t Print tlib version")
	case "version":
	case "-v":
		fmt.Printf("tlib %s\n", Version)
	}
}
