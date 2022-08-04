package roflmeta

import (
	"errors"
	"path/filepath"
	"strings"
)

var errInvalidTemplate = errors.New("invalid template")

func substringStart(input string, start int) string {
	runes := []rune(input)
	if start >= len(runes) {
		return ""
	}
	return string(runes[start:])
}

func substringStartEnd(input string, start int, end int) string {
	runes := []rune(input)
	if start >= len(runes) {
		return ""
	}
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

func isVideo(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".mp4" || ext == ".m4v" || ext == ".mkv" || ext == ".webm" || ext == ".mov" || ext == ".avi" || ext == ".wmv" || ext == ".mpg" || ext == ".flv" || ext == ".3gp" {
		return true
	}
	return false
}
