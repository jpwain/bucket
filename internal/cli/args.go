package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	LeftPath  string
	RightPath string
}

func ParseArgs(argv []string) (Args, error) {
	if len(argv) != 2 {
		return Args{}, errors.New("usage: bucket <left-file> <right-file>")
	}
	left, err := filepath.Abs(argv[0])
	if err != nil {
		return Args{}, fmt.Errorf("resolve left path: %w", err)
	}
	right, err := filepath.Abs(argv[1])
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
