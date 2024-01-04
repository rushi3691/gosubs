package f

import (
	"context"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type Translator struct {
	language string
	apiKey   string
	src      language.Tag
	dst      language.Tag
	client   *translate.Client
}

func NewTranslator(language, apiKey string, src, dst language.Tag) (*Translator, error) {
	ctx := context.Background()
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Translator{
		language: language,
		apiKey:   apiKey,
		src:      src,
		dst:      dst,
		client:   client,
	}, nil
}

func (t *Translator) Translate(sentence string) (string, error) {
	if sentence == "" {
		return "", nil
	}

	ctx := context.Background()
	resp, err := t.client.Translate(ctx, []string{sentence}, t.dst, &translate.Options{Source: t.src, Format: translate.Text})
	if err != nil {
		return "", err
	}

	if len(resp) > 0 {
		return resp[0].Text, nil
	}

	return "", nil
}
