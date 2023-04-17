package service

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"ha-video-translator/pkg/gcp"
	"ha-video-translator/pkg/tesseract"
)

type (
	TranslateRequest struct {
		Paths []string `json:"paths"`
	}

	TranslateResponse struct {
		Translations map[string]string `json:"translations"`
	}
	Service struct {
		translator *gcp.Translator
	}

	translation struct {
		source      string
		translation string
	}
)

func New() *Service {
	return &Service{
		translator: gcp.New(context.Background()),
	}
}

func (s *Service) RegisterHttpHandlers() {
	http.HandleFunc("/translate", s.handle)
}

func (s *Service) translate(source string, out chan<- translation, wg *sync.WaitGroup) {
	defer wg.Done()
	ocr := tesseract.Operation{
		PathToFile:           source,
		TargetLanguage:       "chi_tra",
		LogLevel:             tesseract.ERROR,
		PageSegmentationMode: tesseract.FULLY_AUTO_PAGE_SEGMENTATION_NO_OSD,
		OCREngineMode:        tesseract.DEFAULT,
	}.Perform()

	if ocr.ExitCode != 0 {
		out <- translation{source: source}
		return
	}

	translations := s.translator.Translate(context.Background(), "en", ocr.StdOut)
	out <- translation{
		source:      source,
		translation: translations[0],
	}
}

func (s *Service) handle(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Failed to read request")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var parsed TranslateRequest
	if err = json.Unmarshal(body, &parsed); err != nil {
		log.Printf("Failed to parse request body [%v]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var wg sync.WaitGroup
	wg.Add(len(parsed.Paths))
	results := make(chan translation, len(parsed.Paths))

	for idx := range parsed.Paths {
		go s.translate(parsed.Paths[idx], results, &wg)
	}

	wg.Wait()
	close(results)

	response := TranslateResponse{Translations: make(map[string]string)}

	for item := range results {
		response.Translations[item.source] = item.translation
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to marshal respoinse [%v]", err.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}
