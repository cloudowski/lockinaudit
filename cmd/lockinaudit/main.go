// Command lockinaudit measures cloud vendor lock-in, offline and read-only.
package main

import (
	"os"

	"github.com/cloudowski/lockinaudit/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
