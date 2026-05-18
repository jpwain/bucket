package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Args struct {
	LeftPath         string
	RightPath        string
	VersionRequested bool
}

func ParseArgs(argv []string) (Args, error) {
	fs := flag.NewFlagSet("bucket", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	version := fs.Bool("version", false, "print version")
	vShort := fs.Bool("v", false, "print version")

	if err := fs.Parse(argv); err != nil {
		return Args{}, errors.New("usage: bucket [-v|--version] <left-file> <right-file>")
	}

	if *version || *vShort {
		return Args{VersionRequested: true}, nil
	}

	positional := fs.Args()
	if len(positional) != 2 {
		return Args{}, errors.New("usage: bucket <left-file> <right-file>")
	}

	left, err := filepath.Abs(positional[0])
	if err != nil {
		return Args{}, fmt.Errorf("resolve left path: %w", err)
	}
	right, err := filepath.Abs(positional[1])
	if err != nil {
		return Args{}, fmt.Errorf("resolve right path: %w", err)
	}

	ls, err := os.Stat(left)
	if err != nil {
		return Args{}, fmt.Errorf("stat left file: %w", err)
	}
	rs, err := os.Stat(right)
	if err != nil {
		return Args{}, fmt.Errorf("stat right file: %w", err)
	}
	if os.SameFile(ls, rs) {
		return Args{}, errors.New("left and right paths must refer to different files")
	}

	return Args{LeftPath: left, RightPath: right}, nil
}
