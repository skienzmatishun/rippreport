package main

import (
	"os"

	"github.com/rippreport/go-content-pipeline/internal/cli"
)

func main() {
	app := cli.NewApp(os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args))
}

