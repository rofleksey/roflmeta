package roflmeta

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
)

var errInvalidTemplate = errors.New("invalid template")
var errRegexFailed = errors.New("regex failed")
var errMultipleFailed = errors.New("multiple method failed")

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

func longestCommonPrefix(arr []string) int {
	if len(arr) == 0 {
		return -1
	}
	if len(arr) == 1 {
		return len(arr)
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	firstRunes := []rune(arr[0])
	lastRunes := []rune(arr[len(arr)-1])
	minLen := len(firstRunes)
	if len(lastRunes) < minLen {
		minLen = len(lastRunes)
	}

	i := 0
	for i < minLen && firstRunes[i] == lastRunes[i] {
		i++
	}
	return i
}
