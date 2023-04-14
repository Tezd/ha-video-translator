package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type (
	Translator struct {
		client *translate.Client
	}
)

func (t *Translator) translate(ctx context.Context, targetLanguage string, text ...string) []string {

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		log.Printf("Invalid translate language provided [%v] [%v]", targetLanguage, err)
		return nil
	}

	response, err := t.client.Translate(ctx, text, lang, nil)
	if err != nil {
		log.Printf("Translation failed for text [%v] for language [%v] with error [%v]", text, targetLanguage, err)
		return nil
	}

	if len(response) == 0 {
		log.Printf("Translation returned zero results for text [%v] for language [%v]", text, targetLanguage)
		return nil
	}

	var translated []string

	for idx := range response {
		translated = append(translated, response[idx].Text)
	}

	return translated
}

func New(ctx context.Context) *Translator {

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize translation client [%v]", err)
		return nil
	}

	return &Translator{
		client: client,
	}
}
