package tesseract

import (
	"io"
	"os/exec"
	"strconv"
)

const (
	ERROR_EXIT_CODE = 1
)

type (
	PageSegmentation = uint8
	OCREngineMode    = uint8
	LogLevel         = string
)

const (
	OSD_ONLY PageSegmentation = iota
	AUTO_PAGE_SEGMENTATION_WITH_OSD
	ONLY_AUTO_PAGE_SEGMENTATION
	FULLY_AUTO_PAGE_SEGMENTATION_NO_OSD
	SINGLE_COLUMN_OF_VARIABLE_SIZE
	SINGLE_UNIFORM_BLOCK_OF_VERT_ALIGNED_TEXT
	SINGLE_UNIFORM_BLOCK_OF_TEXTX
	IMAGE_AS_SINGLE_TEXT_LINE
	IMAGE_AS_SINGLE_WORD
	IMAGE_AS_SINGLE_WORD_IN_CIRCLE
	IMAGE_AS_SINGLE_CHARACTER
	UNORDERED_SPARSE_TEXT
	SPARSE_TEXT_WITH_OSD
	RAW_LINE
)

const (
	LEGACY OCREngineMode = iota
	LSTM_ONLY
	LEGACY_AND_LSTM
	DEFAULT
)

const (
	ALL     LogLevel = "ALL"
	TRACE            = "TRACE"
	DEBUG            = "DEBUG"
	INFO             = "INFO"
	WARNING          = "WARN"
	ERROR            = "ERROR"
	FATAL            = "FATAL"
	OFF              = "OFF"
)

type (
	Operation struct {
		PathToFile           string
		TargetLanguage       string
		LogLevel             LogLevel
		PageSegmentationMode PageSegmentation
		OCREngineMode        OCREngineMode
	}

	Result struct {
		StdOut   string
		StdErr   string
		ExitCode int
	}
)

func (o Operation) Perform() Result {
	cmd := exec.Command(
		"tesseract",
		[]string{
			o.PathToFile,
			"-",
			"-l", o.TargetLanguage,
			"--psm", strconv.Itoa(int(o.PageSegmentationMode)),
			"--oem", strconv.Itoa(int(o.OCREngineMode)),
			"--loglevel", o.LogLevel,
		}...,
	)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return makeErrorResult(err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return makeErrorResult(err)
	}

	if err = cmd.Start(); err != nil {
		return makeErrorResult(err)
	}

	outBytes, err := io.ReadAll(stdOut)
	if err != nil {
		return makeErrorResult(err)
	}

	errBytes, err := io.ReadAll(stdErr)
	if err != nil {
		return makeErrorResult(err)
	}

	result := Result{
		StdOut: string(outBytes),
		StdErr: string(errBytes),
	}

	if err = cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitError.ExitCode()
		} else {
			result.ExitCode = ERROR_EXIT_CODE
		}
	}

	return result
}

func makeErrorResult(err error) Result {
	return Result{
		StdErr:   err.Error(),
		ExitCode: ERROR_EXIT_CODE,
	}
}
