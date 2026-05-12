package main

import (
	"fmt"
	"os"

	"bucket/internal/app"
	"bucket/internal/cli"
	fileio "bucket/internal/io"
)

func main() {
	args, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	left, err := fileio.LoadBucket(args.LeftPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	right, err := fileio.LoadBucket(args.RightPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := app.Run(left, right); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
