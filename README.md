# Installation

```
$ go get github.com/noffle/easy-ipfs-shell
```

# Usage
```
import (
  shell "github.com/noffle/easy-ipfs-shell/shell"
)
```

### shell.New() (shell.Shell, error)

Returns a Shell interface, preferring a local HTTP API node if it can find one,
but falling back to producing a new ephemeral node that self-bootstraps.

