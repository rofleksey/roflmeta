package roflmeta

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func genInput(format string, from int, to int) []string {
	filenames := make([]string, 0, to-from+1)
	for i := from; i <= to; i++ {
		filenames = append(filenames, fmt.Sprintf(format, i))
	}
	return filenames
}

func genOutput(season string, format string, from int, to int) []EpisodeMetadata {
	result := make([]EpisodeMetadata, 0, to-from+1)
	for i := from; i <= to; i++ {
		result = append(result, EpisodeMetadata{
			Season:  season,
			Episode: fmt.Sprintf(format, i),
		})
	}
	return result
}

func genSingle(season string, episode string) EpisodeMetadata {
	return EpisodeMetadata{
		Season:  season,
		Episode: episode,
	}
}

func TestMultipleEpisodeMetadata1(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Hunter x Hunter (2011) - Episodes 001-148/[Judas] Hunter x Hunter (2011) - S01E%03d.mkv", 1, 148)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "%03d", 1, 148)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata2(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("Dr Stone/Dr Stone - %02d.mkv", 1, 24)...)
	input = append(input, genInput("Dr Stone Season 2/Dr Stone - %02d.mkv", 1, 11)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("", "%02d", 1, 24)...)
	expected = append(expected, genOutput("Season 2", "%02d", 1, 11)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata3(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Anime Time] Durarara!!/[Anime Time] Durarara!! - %02d.mkv", 1, 12)...)
	input = append(input, "[Anime Time] Durarara!!/[Anime Time] Durarara!! - 12.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!!/[Anime Time] Durarara!! - %02d.mkv", 13, 25)...)
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ketsu/[Anime Time] Durarara!! X2 Ketsu - %02d.mkv", 1, 12)...)
	input = append(input, "[Anime Time] Durarara!! X2 Ketsu/[Anime Time] Durarara!! X2 Ketsu - 7.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Shou/[Anime Time] Durarara!! X2 Shou - %02d.mkv", 1, 12)...)
	input = append(input, "[Anime Time] Durarara!! X2 Shou/[Anime Time] Durarara!! X2 Shou - OVA.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - %02d.mkv", 1, 9)...)
	input = append(input, "[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - 1.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - %02d.mkv", 10, 12)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("", "%02d", 1, 12)...)
	expected = append(expected, genSingle("", "12.5"))
	expected = append(expected, genOutput("", "%02d", 13, 25)...)
	expected = append(expected, genOutput("X2 Ketsu", "%02d", 1, 12)...)
	expected = append(expected, genSingle("X2 Ketsu", "7.5"))
	expected = append(expected, genOutput("X2 Shou", "%02d", 1, 12)...)
	expected = append(expected, genSingle("X2 Shou", "OVA"))
	expected = append(expected, genOutput("X2 Ten", "%02d", 1, 9)...)
	expected = append(expected, genSingle("X2 Ten", "1.5"))
	expected = append(expected, genOutput("X2 Ten", "%02d", 10, 12)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata4(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Jujutsu Kaisen - S01E%02d.mkv", 1, 24)...)
	input = append(input, "[Judas] Jujutsu Kaisen - S01SP1.mkv")

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "E%02d", 1, 24)...)
	expected = append(expected, genSingle("01", "SP1"))

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata5(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("season 01/episode %02d.mkv", 1, 24)...)
	input = append(input, genInput("season 02/episode %02d.mkv", 1, 24)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "%02d", 1, 24)...)
	expected = append(expected, genOutput("2", "%02d", 1, 24)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata6(t *testing.T) {
	input := make([]string, 0, 256)
	for i := 1; i <= 24; i++ {
		input = append(input, fmt.Sprintf("nartsiss %02d/nartsiss %02d.mkv", i, i))
	}

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("nartsiss", "%02d", 1, 24)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata7(t *testing.T) {
	input := make([]string, 0, 256)
	for j := 1; j <= 2; j++ {
		for i := 1; i <= 24; i++ {
			input = append(input, fmt.Sprintf("nartsiss %02d/nartsiss %02d-%02d.mkv", j, i, j))
		}
	}
	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "%02d-01", 1, 24)...)
	expected = append(expected, genOutput("2", "%02d-02", 1, 24)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata8(t *testing.T) {
	input := make([]string, 0, 256)
	for j := 1; j <= 2; j++ {
		for i := 1; i <= 24; i++ {
			input = append(input, fmt.Sprintf("nartsiss %02d/nartsiss %02dx%02d.mkv", j, j, i))
		}
	}
	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "1x%02d", 1, 24)...)
	expected = append(expected, genOutput("2", "2x%02d", 1, 24)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if !success {
		t.Fatal("Used single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata9(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Anime Time] Durarara!!/[Anime Time] Durarara!! - %02d.mkv", 1, 12)...)
	input = append(input, "[Anime Time] Durarara!!/[Anime Time] Durarara!! - 12.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!!/[Anime Time] Durarara!! - %02d.mkv", 13, 25)...)
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ketsu/[Anime Time] Durarara!! X2 Ketsu - %02d.mkv", 1, 12)...)
	// intentional typo in dir here
	input = append(input, "[Anime Tim] Durarara!! X2 Ketsu/[Anime Time] Durarara!! X2 Ketsu - 7.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Shou/[Anime Time] Durarara!! X2 Shou - %02d.mkv", 1, 12)...)
	input = append(input, "[Anime Time] Durarara!! X2 Shou/[Anime Time] Durarara!! X2 Shou - OVA.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - %02d.mkv", 1, 9)...)
	input = append(input, "[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - 1.5.mkv")
	input = append(input, genInput("[Anime Time] Durarara!! X2 Ten/[Anime Time] Durarara!! X2 Ten - %02d.mkv", 10, 12)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("Durarara!!", "%02d", 1, 12)...)
	expected = append(expected, genSingle("Durarara!!", "12.5"))
	expected = append(expected, genOutput("Durarara!!", "%02d", 13, 25)...)
	expected = append(expected, genOutput("Durarara!! X2 Ketsu", "%02d", 1, 12)...)
	expected = append(expected, genSingle("Durarara!! X2 Ketsu", "7.5"))
	expected = append(expected, genOutput("Durarara!! X2 Shou", "%02d", 1, 12)...)
	expected = append(expected, genSingle("Durarara!! X2 Shou", "OVA"))
	expected = append(expected, genOutput("Durarara!! X2 Ten", "%02d", 1, 9)...)
	expected = append(expected, genSingle("Durarara!! X2 Ten", "1.5"))
	expected = append(expected, genOutput("Durarara!! X2 Ten", "%02d", 10, 12)...)

	metadataArr, success := ParseMultipleEpisodeMetadata(input)
	if success {
		t.Fatal("Expected single parser here")
	}
	assert.Equal(t, metadataArr, expected)
}
