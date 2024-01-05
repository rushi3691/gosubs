package u

import (
	"time"

	"github.com/asticode/go-astisub"
)

const (
	DefaultPaddingAfter  = 0.0
	DefaultPaddingBefore = 0.0
)

//

func SrtFormatter(subtitles []Subtitle, paddingBefore float64, paddingAfter float64, outputFilePath string) error {
	srtSubtitles := astisub.NewSubtitles()
	for i, subtitle := range subtitles {
		start := subtitle.Region[0] - paddingBefore // in seconds
		if start < 0 {
			start = 0
		}
		end := subtitle.Region[1] + paddingAfter // in seconds

		item := &astisub.Item{
			Index:   i + 1,
			StartAt: time.Duration(start * float64(time.Second)),
			EndAt:   time.Duration(end * float64(time.Second)),
			Lines: []astisub.Line{
				{
					Items: []astisub.LineItem{
						{
							Text: subtitle.Transcript,
						},
					},
				},
			},
		}
		srtSubtitles.Items = append(srtSubtitles.Items, item)
	}

	err := srtSubtitles.Write(outputFilePath)
	if err != nil {
		return err
	}

	return nil
}
