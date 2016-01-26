package shell

import (
	"fmt"
	"os"
	"os/signal"

	api "github.com/ipfs/go-ipfs-api"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	embedded "github.com/whyrusleeping/ipfs-embedded-shell"
	context "golang.org/x/net/context"
)

func NewShell() (Shell, error) {
	myShell, err := getApiShell()
	if err == nil {
		// fmt.Println("got an api shell!")
		return myShell, nil
	}

	myShell, err = getEmbeddedShell()
	if err == nil {
		// fmt.Println("got an embedded shell!")
		return myShell, nil
	}

	return nil, err
}

func getApiShell() (Shell, error) {
	apiShell := api.NewShell("http://127.0.0.1:5001")
	_, _, err := apiShell.Version()
	if err != nil {
		return nil, err
	}

	return apiShell, nil
}

func getEmbeddedShell() (Shell, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the ipfs node context if the process gets interrupted or killed.
	// TODO(noffle): is this needed?
	go func() {
		interrupts := make(chan os.Signal, 1)
		signal.Notify(interrupts, os.Interrupt, os.Kill)
		<-interrupts
		cancel()
	}()

	repoPath, err := getRepoPath()
	if err != nil {
		return nil, fmt.Errorf("couldn't get repo path: %s", err)
	}

	embeddedShell, err := embedded.NewDefaultNodeWithFSRepo(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't get embedded shell: %s", err)
	}

	return embedded.NewShell(embeddedShell), nil
}

func getRepoPath() (string, error) {
	repoPath, err := fsrepo.BestKnownPath()
	if err != nil {
		return "", err
	}
	return repoPath, nil
}
