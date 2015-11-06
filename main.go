package main

import (
	"fmt"
	"io"
	"os"

	shell "github.com/noffle/easy-ipfs-shell/shell"
)

func main() {
	shell, err := shell.NewShell()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("got a generic shell!")

	reader, err := shell.Cat("QmVVjWrps58cFS1hSvCdAxmS4wggKfRGbDzJway6QCxR4U")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	io.Copy(os.Stdout, reader)
}
