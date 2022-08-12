package roflmeta

import (
	"regexp"
	"sort"
	"strings"
)

var spaceRegex = regexp.MustCompile("^\\s*$")
var bracketRemoveRegex = regexp.MustCompile("\\[.*?\\]|\\(.*?\\)|\\{.*?\\}")

type frequency struct {
	value int
	group int
}

func preCleanFileName(filename string) string {
	return bracketRemoveRegex.ReplaceAllLiteralString(filename, "")
}

func postCleanData(data string) string {
	data = strings.Trim(data, " -_/\\*.'")
	data = strings.TrimSpace(data)
	return data
}

func cleanFileNames(filename []string) []string {
	result := make([]string, 0, len(filename))
	for _, name := range filename {
		result = append(result, preCleanFileName(name))
	}
	return result
}

// calculates number of distinct matches for each group and sorts them
func calcRegexFrequencies(filenames []string, regex *regexp.Regexp, groupCount int) []frequency {
	result := make([]frequency, 0, groupCount)
	for group := 1; group <= groupCount; group++ {
		set := make(map[string]struct{}, len(filenames))
		for _, name := range filenames {
			test := regex.FindStringSubmatch(name)
			// don't count empty matches
			if !spaceRegex.MatchString(test[group]) {
				set[test[group]] = struct{}{}
			}
		}
		f := frequency{
			value: len(set),
			group: group,
		}
		result = append(result, f)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].value < result[j].value
	})
	return result
}

func calcDistinctFrequencies(frequencies []frequency) int {
	set := make(map[int]struct{}, len(frequencies))
	for _, f := range frequencies {
		set[f.value] = struct{}{}
	}
	return len(set)
}

func testFreqGroupMonotonous(frequencies []frequency) bool {
	for i := 1; i < len(frequencies); i++ {
		if frequencies[i-1].group >= frequencies[i].group {
			return false
		}
	}
	return true
}

func parseChangingEpisodes(filenames []string, testSeasonFilename string, regex *regexp.Regexp, episodeGroup int) []EpisodeMetadata {
	result := make([]EpisodeMetadata, 0, len(filenames))
	// will trust try-hard single episode parser on this one
	season := ParseSingleEpisodeMetadata(testSeasonFilename).Season
	for _, name := range filenames {
		test := regex.FindStringSubmatch(name)
		result = append(result, EpisodeMetadata{
			Episode: postCleanData(test[episodeGroup]),
			Season:  season,
		})
	}
	return result
}

func parseEpisodesAndSeasons(filenames []string, regex *regexp.Regexp, seasonGroup int, episodeGroup int) []EpisodeMetadata {
	result := make([]EpisodeMetadata, 0, len(filenames))
	for _, name := range filenames {
		test := regex.FindStringSubmatch(name)
		result = append(result, EpisodeMetadata{
			Episode: postCleanData(test[episodeGroup]),
			Season:  postCleanData(test[seasonGroup]),
		})
	}
	return result
}

func fallbackToSingleParser(filenames []string) []EpisodeMetadata {
	result := make([]EpisodeMetadata, 0, len(filenames))
	for _, name := range filenames {
		result = append(result, ParseSingleEpisodeMetadata(name))
	}
	return result
}

func restoreFullAnswer(allFilenames []string, videoFilenames []string, partialAnswer []EpisodeMetadata) []EpisodeMetadata {
	fullResult := make([]EpisodeMetadata, 0, len(allFilenames))
	videoIndex := 0
	for _, name := range allFilenames {
		if videoIndex < len(videoFilenames) && name == videoFilenames[videoIndex] {
			fullResult = append(fullResult, partialAnswer[videoIndex])
			videoIndex++
		} else {
			fullResult = append(fullResult, EpisodeMetadata{})
		}
	}
	return fullResult
}

// ParseMultipleEpisodeMetadata attempts to parse metadata from multiple filenames
// See EpisodeMetadata for details
// It tries to figure out filenames' template and gather information according to it
// Returned bool indicates whether this attempt was successful (true), or single episode parser was used as a fallback (false)
func ParseMultipleEpisodeMetadata(filenames []string) ([]EpisodeMetadata, bool) {
	if len(filenames) == 0 {
		return []EpisodeMetadata{}, true
	}
	if len(filenames) == 1 {
		return []EpisodeMetadata{ParseSingleEpisodeMetadata(filenames[0])}, true
	}
	cleanedFilenames := cleanFileNames(filenames)
	videoFilenames := make([]string, 0, len(cleanedFilenames))

	for _, name := range cleanedFilenames {
		if isVideo(name) {
			videoFilenames = append(videoFilenames, name)
		}
	}

	t, err := restoreTemplate(videoFilenames)

	if err != nil {
		// something went wrong, fallback to single
		return fallbackToSingleParser(filenames), false
	}

	varCount := t.varCount()
	regex := t.toRegex()

	// definitely a single season with changing episodes
	if varCount == 1 {
		partialAnswer := parseChangingEpisodes(videoFilenames, filenames[0], regex, 1)
		return restoreFullAnswer(cleanedFilenames, videoFilenames, partialAnswer), true
	}

	frequencies := calcRegexFrequencies(videoFilenames, regex, varCount)
	distinctFreqCount := calcDistinctFrequencies(frequencies)
	groupMonotonous := testFreqGroupMonotonous(frequencies)

	// may be something like /%season/%season/%episode-%season.mkv, try to swap last two groups
	if !groupMonotonous {
		frLen := len(frequencies)
		frequencies[frLen-1], frequencies[frLen-2] = frequencies[frLen-2], frequencies[frLen-1]
		groupMonotonous = testFreqGroupMonotonous(frequencies)
	}

	if groupMonotonous {
		// definitely only episodes
		if distinctFreqCount == 1 {
			partialAnswer := parseChangingEpisodes(videoFilenames, filenames[0], regex, frequencies[len(frequencies)-1].group)
			return restoreFullAnswer(cleanedFilenames, videoFilenames, partialAnswer), true
		}
		// probably seasons and episodes
		if distinctFreqCount == 2 {
			partialAnswer := parseEpisodesAndSeasons(videoFilenames, regex, frequencies[len(frequencies)-2].group, frequencies[len(frequencies)-1].group)
			return restoreFullAnswer(cleanedFilenames, videoFilenames, partialAnswer), true
		}
	}

	// fallback to good ol' single parser
	return fallbackToSingleParser(filenames), false
}
