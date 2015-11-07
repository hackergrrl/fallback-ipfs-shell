package shell

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestCat(t *testing.T) {
	myShell, err := NewShell()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	reader, err := myShell.Cat("QmQLBvJ3ur7U7mzbYDLid7WkaciY84SLpPYpGPHhDNps2Y")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	output := buf.String()

	expected := "\"The man who makes no mistakes does not make anything.\" - Edward John Phelps\n"

	if output != expected {
		t.FailNow()
	}
}
