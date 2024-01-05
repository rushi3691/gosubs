package services

import (
	"os"
	"os/exec"
	"strconv"
)

// include_before=0.25, include_after=0.25

const (
	DefaultIncludeBefore = 0.25
	DefaultIncludeAfter  = 0.25
)

type FLACConverter struct {
	sourcePath    string
	includeBefore float64
	includeAfter  float64
}

func NewFLACConverter(sourcePath string, includeBefore, includeAfter float64) *FLACConverter {
	return &FLACConverter{
		sourcePath:    sourcePath,
		includeBefore: includeBefore,
		includeAfter:  includeAfter,
	}
}

func (f *FLACConverter) Convert(region [2]float64) ([]byte, error) {
	start, end := region[0], region[1]
	start = max(0, start-f.includeBefore)
	end += f.includeAfter

	tempFile, err := os.CreateTemp("", "*.flac")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	command := []string{"ffmpeg", "-ss", strconv.FormatFloat(start, 'f', -1, 64), "-t", strconv.FormatFloat(end-start, 'f', -1, 64),
		"-y", "-i", f.sourcePath,
		"-loglevel", "error", tempFile.Name()}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = nil

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	readData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, err
	}

	return readData, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
