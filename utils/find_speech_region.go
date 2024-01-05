package utils

import (
	"log"
	"math"
	"os"

	"github.com/faiface/beep/wav"
	"github.com/montanaflynn/stats"
)

// frame_width=4096, min_region_size=0.5, max_region_size=6
const (
	DefaultFrameWidth    = 4096
	DefaultMinRegionSize = 0.5
	DefaultMaxRegionSize = 6
)

func FindSpeechRegions(filename string, frameWidth int, minRegionSize float64, maxRegionSize float64) ([][2]float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer streamer.Close()

	sampleRate := format.SampleRate
	chunkDuration := float64(frameWidth) / float64(sampleRate)
	// sampleWidth := format.Width()
	// channels := format.NumChannels
	// fmt.Println("chunkDuration", chunkDuration)
	// fmt.Println("sampleRate", sampleRate)
	// fmt.Println("sampleWidth", sampleWidth)
	// fmt.Println("channels", channels)

	energies := []float64{}

	buf := make([][2]float64, frameWidth)

	for {
		_, ok := streamer.Stream(buf)
		if !ok {
			break
		}
		energies = append(energies, rms(buf))
	}

	// fmt.Println("energies", energies)
	// fmt.Println("len(energies)", len(energies))

	threshold, _ := stats.Percentile(energies, 20)

	elapsedTime := 0.0
	var regions [][2]float64
	var regionStart float64
	regionStarted := false

	for _, energy := range energies {
		isSilence := energy <= threshold
		maxExceeded := regionStarted && elapsedTime-regionStart >= maxRegionSize

		if (maxExceeded || isSilence) && regionStarted {
			if elapsedTime-regionStart >= minRegionSize {
				regions = append(regions, [2]float64{regionStart, elapsedTime})
				regionStarted = false
			}
		} else if !regionStarted && !isSilence {
			regionStart = elapsedTime
			regionStarted = true
		}
		elapsedTime += chunkDuration
	}
	// return regions
	// fmt.Println(regions)
	// fmt.Println("len(regions)", len(regions))
	return regions, nil
}

// func rms(buf [][2]float64) float64 {
// 	sum := 0.0
// 	for _, v := range buf {
// 		sum += math.Pow(v[0], 2) + math.Pow(v[1], 2)
// 	}
// 	return math.Sqrt(sum/float64(len(buf))) * 46300.0
// }

// func rms(buf [][2]float64) float64 {
// 	sum := 0.0
// 	for _, v := range buf {
// 		scaledV0 := v[0] * 32767.0
// 		scaledV1 := v[1] * 32767.0
// 		sum += math.Pow(scaledV0, 2) + math.Pow(scaledV1, 2)
// 	}
// 	return math.Sqrt(sum / float64(len(buf)))
// }

func rms(buf [][2]float64) float64 {
	sum := 0.0
	for _, v := range buf {
		// Convert float64 values between -1 and 1 back to int16
		left := int16(v[0] * 32767)
		right := int16(v[1] * 32767)
		sum += math.Pow(float64(left), 2) + math.Pow(float64(right), 2)
	}
	return math.Sqrt(sum/float64(len(buf)*2)) * 2 // Multiply by 2 because it's stereo data
}
