package main

import (
	"context"

	"github.com/iFleksy/gopassxc/cmd/root"
)

func main() {
	root.GetRootCMD().ExecuteContext(context.Background())
}
