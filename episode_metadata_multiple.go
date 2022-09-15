package roflmeta

import (
	"path/filepath"
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

type fileEntry struct {
	cleanedFileName string
	dir             string
	isVideo         bool
	result          EpisodeMetadata
}

func preCleanFileName(filename string) string {
	return bracketRemoveRegex.ReplaceAllLiteralString(filename, "")
}

func postCleanData(data string) string {
	data = strings.Trim(data, " -_/\\*.'")
	data = strings.TrimSpace(data)
	return data
}

// calculates number of distinct matches for each group and sorts them
func calcRegexFrequencies(filenames []string, regex *regexp.Regexp, groupCount int) ([]frequency, error) {
	result := make([]frequency, 0, groupCount)
	for group := 1; group <= groupCount; group++ {
		set := make(map[string]struct{}, len(filenames))
		for _, name := range filenames {
			test := regex.FindStringSubmatch(name)
			if test == nil {
				return nil, errRegexFailed
			}
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
	return result, nil
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

func parseChangingDirs(dirs []string, regex *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	for _, dirName := range dirs {
		test := regex.FindStringSubmatch(dirName)
		result[dirName] = postCleanData(test[1])
	}
	return result
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

func getSeasonSet(dirFileMap map[string][]*fileEntry) map[string]struct{} {
	seasonSet := make(map[string]struct{})
	for _, entries := range dirFileMap {
		for _, entry := range entries {
			seasonSet[entry.result.Season] = struct{}{}
		}
	}
	return seasonSet
}

func getSeasons(dirFileMap map[string][]*fileEntry) []string {
	seasons := make([]string, 0, len(dirFileMap))
	seasonSet := getSeasonSet(dirFileMap)
	for season := range seasonSet {
		seasons = append(seasons, season)
	}
	return seasons
}

func getDirs(dirFileMap map[string][]*fileEntry) []string {
	dirs := make([]string, 0, len(dirFileMap))
	for dir, _ := range dirFileMap {
		dirs = append(dirs, dir)
	}
	return dirs
}

func parseMultipleEpisodeMetadataImpl(filenames []string) ([]EpisodeMetadata, error) {
	if len(filenames) == 1 {
		return fallbackToSingleParser(filenames), nil
	}

	t, err := restoreTemplate(filenames)

	if err != nil {
		return nil, err
	}

	varCount := t.varCount()
	regex := t.toRegex()

	// definitely a single season with changing episodes
	if varCount == 1 {
		return parseChangingEpisodes(filenames, filenames[0], regex, 1), nil
	}

	frequencies, err := calcRegexFrequencies(filenames, regex, varCount)
	if err != nil {
		return nil, err
	}

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
			return parseChangingEpisodes(filenames, filenames[0], regex, frequencies[len(frequencies)-1].group), nil
		}
		// probably seasons and episodes
		if distinctFreqCount == 2 {
			return parseEpisodesAndSeasons(filenames, regex, frequencies[len(frequencies)-2].group, frequencies[len(frequencies)-1].group), nil
		}
	}

	return nil, errMultipleFailed
}

// ParseMultipleEpisodeMetadata attempts to parse metadata from multiple filenames
// See EpisodeMetadata for details
// It tries to figure out filenames' template and gather information according to it
func ParseMultipleEpisodeMetadata(filenames []string) []EpisodeMetadata {
	if len(filenames) == 0 {
		return []EpisodeMetadata{}
	}
	if len(filenames) == 1 {
		return []EpisodeMetadata{ParseSingleEpisodeMetadata(filenames[0])}
	}

	// process files in each dir separately
	fileEntries := make([]*fileEntry, 0, len(filenames))
	dirFileMap := make(map[string][]*fileEntry)
	for _, name := range filenames {
		entry := &fileEntry{
			cleanedFileName: preCleanFileName(name),
			dir:             filepath.Dir(name),
			isVideo:         isVideo(name),
		}
		fileEntries = append(fileEntries, entry)
		if entry.isVideo {
			list := dirFileMap[entry.dir]
			list = append(list, entry)
			dirFileMap[entry.dir] = list
		}
	}

	for _, entries := range dirFileMap {
		dirFilenames := make([]string, 0, len(filenames))
		for _, f := range entries {
			dirFilenames = append(dirFilenames, f.cleanedFileName)
		}
		dirResult, err := parseMultipleEpisodeMetadataImpl(dirFilenames)
		if err != nil {
			dirResult = fallbackToSingleParser(dirFilenames)
		}
		for i, r := range dirResult {
			entries[i].result = r
		}
	}

	dirs := getDirs(dirFileMap)
	if len(dirs) > 1 {
		seasonSet := getSeasonSet(dirFileMap)
		// multiple dirs AND single season, decide by dirname
		if len(seasonSet) == 1 {
			t, err := restoreTemplate(dirs)
			if err == nil && t.varCount() == 1 {
				seasonsMap := parseChangingDirs(dirs, t.toRegex())
				for dir, entries := range dirFileMap {
					for _, entry := range entries {
						entry.result.Season = seasonsMap[dir]
					}
				}
			}
		}
	}

	// remove seasonal common prefix
	seasons := getSeasons(dirFileMap)
	lcp := longestCommonPrefix(seasons)
	if len(seasons) > 1 && lcp > 0 {
		for _, entries := range dirFileMap {
			for _, entry := range entries {
				// lcp < len
				if lcp <= len(entry.result.Season) {
					entry.result.Season = postCleanData(substringStart(entry.result.Season, lcp))
				}
			}
		}
	}

	result := make([]EpisodeMetadata, 0, len(filenames))
	for _, entry := range fileEntries {
		result = append(result, entry.result)
	}
	return result
}
