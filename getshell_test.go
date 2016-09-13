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

	reader, err := myShell.Cat("QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	output := buf.String()

	expected := `Come hang out in our IRC chat room if you have any questions.

Contact the ipfs dev team:
- Bugs: https://github.com/ipfs/go-ipfs/issues
- Help: irc.freenode.org/#ipfs
- Email: dev@ipfs.io
`

	if output != expected {
		t.FailNow()
	}
}
