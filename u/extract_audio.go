package u

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	// DefaultChannels is the default number of channels for the audio
	DefaultChannels = 1
	// DefaultRate is the default sample rate for the audio
	DefaultRate = 16000
)

func ExtractAudio(filename string, channels int, rate int) (string, int, error) {
	tempFile, err := os.CreateTemp("", "*.wav")
	if err != nil {
		return "", 0, err
	}
	defer tempFile.Close()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", 0, fmt.Errorf("the given file does not exist: %s", filename)
	}

	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		return "", 0, fmt.Errorf("ffmpeg: executable not found on machine")
	}

	command := []string{"ffmpeg", "-y", "-i", filename,
		"-ac", fmt.Sprint(channels), "-ar", fmt.Sprint(rate),
		"-loglevel", "error", tempFile.Name()}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", filepath.Join(command...))
	} else {
		cmd = exec.Command(command[0], command[1:]...)
	}

	_, err = cmd.Output()
	if err != nil {
		return "", 0, err
	}

	return tempFile.Name(), rate, nil
}
