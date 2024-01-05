package utils

import (
	"testing"
)

func TestExtractAudio(t *testing.T) {
	_, _, err := ExtractAudio("nonexistent", DefaultChannels, DefaultRate)
	if err == nil {
		t.Errorf("ExtractAudio() should have returned an error")
	}
}

func TestExtractAudio2(t *testing.T) {
	filename, rate, err := ExtractAudio("../sample.mp4", DefaultChannels, DefaultRate)
	if err != nil {
		t.Errorf("ExtractAudio() should not have returned an error")
	}
	if filename == "" {
		t.Errorf("ExtractAudio() should have returned a filename")
	}
	if rate != DefaultRate {
		t.Errorf("ExtractAudio() should have returned a rate of %d", DefaultRate)
	}

	t.Log(filename)
}
