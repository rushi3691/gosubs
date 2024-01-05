package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// language="en", rate=44100, retries=3, api_key=GOOGLE_SPEECH_API_KEY

const (
	GOOGLE_SPEECH_API_URL = "http://www.google.com/speech-api/v2/recognize?client=chromium&lang=%s&key=%s"
	// GoogleSpeechAPIURL = "https://speech.googleapis.com/v1/speech:recognize?lang=%s&key=%s"
	DefaultRate     = 44100
	DefaultRetries  = 3
	DefaultLanguage = "en"
)

type SpeechRecognizer struct {
	Language string
	Rate     int
	APIKey   string
	Retries  int
}

type SpeechRecognitionResult struct {
	Result []struct {
		Alternative []struct {
			Transcript string `json:"transcript"`
		} `json:"alternative"`
	} `json:"result"`
}

func NewSpeechRecognizer(language string, rate, retries int, apiKey string) *SpeechRecognizer {
	return &SpeechRecognizer{
		Language: language,
		Rate:     rate,
		APIKey:   apiKey,
		Retries:  retries,
	}
}

func makeHttpPost(url string, headers map[string]string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func (s *SpeechRecognizer) Recognize(data []byte) (string, error) {
	for i := 0; i < s.Retries; i++ {
		url := fmt.Sprintf(GOOGLE_SPEECH_API_URL, s.Language, s.APIKey)
		headers := map[string]string{
			"Content-Type": fmt.Sprintf("audio/x-flac; rate=%d", s.Rate),
		}

		resp, err := makeHttpPost(url, headers, data)
		if err != nil {
			log.Println(err)
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		// log.Println(string(body))
		lines := bytes.Split(body, []byte("\n"))
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}

			var result SpeechRecognitionResult
			err := json.Unmarshal(line, &result)
			if err != nil {
				log.Println(err)
				continue
			}
			if len(result.Result) > 0 && len(result.Result[0].Alternative) > 0 {
				transcript := result.Result[0].Alternative[0].Transcript
				mod := strings.ToUpper(transcript[:1]) + transcript[1:]
				// log.Println(mod)
				return mod, nil
			}
		}

		return "", nil
	}
	return "", fmt.Errorf("failed to recognize speech")
}
