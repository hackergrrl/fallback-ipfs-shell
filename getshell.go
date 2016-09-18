package shell

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	api "github.com/ipfs/go-ipfs-api"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	embedded "github.com/noffle/ipfs-embedded-shell"
	context "golang.org/x/net/context"
)

func NewShell() (Shell, error) {
	myShell, err := NewApiShell()
	if err == nil {
		// fmt.Println("got an api shell!")
		return myShell, nil
	}

	myShell, err = NewEmbeddedShell()
	if err == nil {
		// fmt.Println("got an embedded shell!")
		return myShell, nil
	}

	return nil, err
}

func NewApiShell() (Shell, error) {
	api, err := apiAddr()
	if err != nil {
		return nil, err
	}

	apiShell := api.NewShell(api)
	_, _, err := apiShell.Version()
	if err != nil {
		return nil, err
	}

	return apiShell, nil
}

func NewEmbeddedShell() (Shell, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the ipfs node context if the process gets interrupted or killed.
	// TODO(noffle): is this needed?
	go func() {
		interrupts := make(chan os.Signal, 1)
		signal.Notify(interrupts, os.Interrupt, os.Kill)
		<-interrupts
		cancel()
	}()

	shell, err := tryLocal(ctx)
	if err == nil {
		return shell, nil
	}

	node, err := embedded.NewTmpDirNode(ctx)
	if err != nil {
		return nil, err
	}

	return embedded.NewShell(node), nil
}

func tryLocal(ctx context.Context) (Shell, error) {
	repoPath, err := getRepoPath()
	if err != nil {
		return nil, err
	}

	node, err := embedded.NewDefaultNodeWithFSRepo(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't get embedded shell: %s", err)
	}

	return embedded.NewShell(node), nil
}

func getRepoPath() (string, error) {
	repoPath, err := fsrepo.BestKnownPath()
	if err != nil {
		return "", err
	}
	return repoPath, nil
}

// APIAddr returns the registered API addr, according to the api file
// in the fsrepo. This is a concurrent operation, meaning that any
// process may read this file. modifying this file, therefore, should
// use "mv" to replace the whole file and avoid interleaved read/writes.
func apiAddr() (string, error) {
	repoPath, err := getRepoPath()
	if err != nil {
		return nil, err
	}

	repoPath = filepath.Clean(repoPath)
	apiFilePath := filepath.Join(repoPath, "api")

	// if there is no file, assume there is no api addr.
	f, err := os.Open(apiFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", repo.ErrApiNotRunning
		}
		return "", err
	}
	defer f.Close()

	// read up to 2048 bytes. io.ReadAll is a vulnerability, as
	// someone could hose the process by putting a massive file there.
	buf := make([]byte, 2048)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	s := string(buf[:n])
	s = strings.TrimSpace(s)
	return s, nil
}
