package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rushi3691/go_subtitle_generator/services"
	"github.com/rushi3691/go_subtitle_generator/utils"
	"golang.org/x/text/language"
)

const (
	DefaultSubtitleFormat = "srt"
)

var (
	DefaultSrcLanguage = language.English
	DefaultDstLanguage = language.English
)

func GenerateSubtitles(
	sourcePath string,
	output_dst string,
	concurrency int,
	srcLanguage string,
	dstLanguage string,
	subtitleFileFormat string,
	apiKey string,
) (string, error) {
	audioFilename, audioRate, err := utils.ExtractAudio(sourcePath, utils.DefaultChannels, utils.DefaultRate)
	if err != nil {
		return "", err
	}
	defer os.Remove(audioFilename)

	regions, err := utils.FindSpeechRegions(audioFilename, utils.DefaultFrameWidth, utils.DefaultMinRegionSize, utils.DefaultMaxRegionSize)
	if err != nil {
		return "", err
	}

	converter := services.NewFLACConverter(audioFilename, services.DefaultIncludeAfter, services.DefaultIncludeBefore)
	// fmt.Println(converter)
	recognizer := services.NewSpeechRecognizer(srcLanguage, audioRate, services.DefaultRetries, apiKey)
	// fmt.Println(recognizer)

	var wg sync.WaitGroup
	subtitles := make([]utils.Subtitle, len(regions))
	regionsChan := make(chan utils.RegionWithIndex)

	// Start 3 worker goroutines
	log.Println("Starting", concurrency, "workers")
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for region := range regionsChan {
				// log.Println("Processing region", region.Index, "in worker", i)
				flacFile, err := converter.Convert(region.Region)
				if err != nil {
					log.Println(err)
					return
				}

				transcript, err := recognizer.Recognize(flacFile)
				if err != nil {
					log.Println(err)
					return
				}

				subtitles[region.Index] = utils.Subtitle{Region: region.Region, Transcript: transcript}
			}
		}(i)
	}

	// Send regions to be processed
	for i, region := range regions {
		regionsChan <- utils.RegionWithIndex{Region: region, Index: i}
	}
	close(regionsChan)

	// Wait for workers to finish
	wg.Wait()

	// Translate subtitles if necessary
	// log.Println(srcLanguage, dstLanguage)
	if srcLanguage != dstLanguage {
		if apiKey != "" {
			translator, err := services.NewTranslator(dstLanguage, apiKey, DefaultSrcLanguage, DefaultDstLanguage)
			if err != nil {
				return "", err
			}

			for i, subtitle := range subtitles {
				translatedTranscript, err := translator.Translate(subtitle.Transcript)
				if err != nil {
					return "", err
				}
				subtitles[i].Transcript = translatedTranscript
			}
		} else {
			return "", fmt.Errorf("subtitle translation requires specified Google Translate API key")
		}
	}

	// formatter := NewFormatter(subtitleFileFormat)
	// formattedSubtitles, err := formatter.Format(subtitles)

	dest := output_dst
	if dest == "" {
		base := strings.TrimSuffix(sourcePath, filepath.Ext(sourcePath))
		dest = fmt.Sprintf("%s.%s", base, subtitleFileFormat)
	}
	log.Println(dest)

	err = utils.SrtFormatter(subtitles, utils.DefaultPaddingBefore, utils.DefaultPaddingAfter, dest)
	if err != nil {
		return "", err
	}

	return dest, nil
}
