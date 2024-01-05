package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_SUBTITLE_FORMAT = "srt"
	DEFAULT_SRC_LANGUAGE    = "en"
	DEFAULT_DST_LANGUAGE    = "en"
)

var (
	DEFAULT_CONCURRENCY = runtime.NumCPU() * 4
	sourcePath          = flag.String("source_path", "", "Path to the video or audio file to subtitle")
	concurrency         = flag.Int("concurrency", DEFAULT_CONCURRENCY, "Number of concurrent API requests to make")
	output              = flag.String("output", "", "Output path for subtitles (by default, subtitles are saved in the same directory and name as the source path)")
	format              = flag.String("format", DEFAULT_SUBTITLE_FORMAT, "Destination subtitle format")
	srcLanguage         = flag.String("src_language", DEFAULT_SRC_LANGUAGE, "Language spoken in source file")
	dstLanguage         = flag.String("dst_language", DEFAULT_DST_LANGUAGE, "Desired language for the subtitles")
	// apiKey      = flag.String("api_key", "", "The Google Translate API key to be used. (Required for subtitle translation)")
)

func ValidateArgs() error {
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	flag.Parse()

	apiKey := os.Getenv("API_KEY")

	subtitleFilePath, err := GenerateSubtitles(*sourcePath, *output, *concurrency, *srcLanguage, *dstLanguage, *format, apiKey)
	if err != nil {
		fmt.Println("Error generating subtitles:", err)
		os.Exit(1)
	}

	fmt.Println("Subtitles file created at", subtitleFilePath)
	os.Exit(0)
}
