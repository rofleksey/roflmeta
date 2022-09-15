package roflmeta

import (
	"fmt"
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

func assertDiff(t *testing.T, actual []EpisodeMetadata, expected []EpisodeMetadata) {
	if len(expected) != len(actual) {
		t.Fatalf("Invalid result length, expected %d, got %d", len(expected), len(actual))
	}
	for i := range expected {
		if expected[i].Season != actual[i].Season {
			t.Fatalf("Invalid season at #%d, expected '%s', got '%s'", i, expected[i].Season, actual[i].Season)
		}
		if expected[i].Episode != actual[i].Episode {
			t.Fatalf("Invalid episode at #%d, expected '%s', got '%s'", i, expected[i].Episode, actual[i].Episode)
		}
	}
}

func TestMultipleEpisodeMetadata1(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Hunter x Hunter (2011) - Episodes 001-148/[Judas] Hunter x Hunter (2011) - S01E%03d.mkv", 1, 148)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "%03d", 1, 148)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata2(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("Dr Stone/Dr Stone - %02d.mkv", 1, 24)...)
	input = append(input, genInput("Dr Stone Season 2/Dr Stone - %02d.mkv", 1, 11)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("", "%02d", 1, 24)...)
	expected = append(expected, genOutput("Season 2", "%02d", 1, 11)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
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

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata4(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Jujutsu Kaisen - S01E%02d.mkv", 1, 24)...)
	input = append(input, "[Judas] Jujutsu Kaisen - S01SP1.mkv")

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "E%02d", 1, 24)...)
	expected = append(expected, genSingle("01", "SP1"))

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata5(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("season 01/episode %02d.mkv", 1, 24)...)
	input = append(input, genInput("season 02/episode %02d.mkv", 1, 24)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "%02d", 1, 24)...)
	expected = append(expected, genOutput("2", "%02d", 1, 24)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata6(t *testing.T) {
	input := make([]string, 0, 256)
	for j := 1; j <= 2; j++ {
		for i := 1; i <= 24; i++ {
			input = append(input, fmt.Sprintf("nartsiss %02d/nartsiss %02d-%02d.mkv", j, i, j))
		}
	}
	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "%02d", 1, 24)...)
	expected = append(expected, genOutput("2", "%02d", 1, 24)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata7(t *testing.T) {
	input := make([]string, 0, 256)
	for j := 1; j <= 2; j++ {
		for i := 1; i <= 24; i++ {
			input = append(input, fmt.Sprintf("nartsiss %02d/nartsiss %02dx%02d.mkv", j, j, i))
		}
	}
	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("1", "%02d", 1, 24)...)
	expected = append(expected, genOutput("2", "%02d", 1, 24)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata8(t *testing.T) {
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

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata9(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("Dr Stone Season 2/Dr Stone - %02d.mkv", 1, 11)...)
	input = append(input, genInput("Dr Stone/Dr Stone - %02d.mkv", 1, 24)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("Season 2", "%02d", 1, 11)...)
	expected = append(expected, genOutput("", "%02d", 1, 24)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata10(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Jujutsu Kaisen - S01E%02d.mkv", 1, 12)...)
	input = append(input, "wrong template 1.txt")
	input = append(input, genInput("[Judas] Jujutsu Kaisen - S01E%02d.mkv", 13, 24)...)
	input = append(input, "[Judas] Jujutsu Kaisen - S01SP1.mkv")
	input = append(input, "[Judas] Jujutsu Kaisen - S01SP1.txt")

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "E%02d", 1, 12)...)
	expected = append(expected, genSingle("", ""))
	expected = append(expected, genOutput("01", "E%02d", 13, 24)...)
	expected = append(expected, genSingle("01", "SP1"))
	expected = append(expected, genSingle("", ""))

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata11(t *testing.T) {
	input := make([]string, 0, 12)
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 01 (1080p) [28B342E5].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 02 (1080p) [2FB81205].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 03 (1080p) [53E3F3D1].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 04 (1080p) [99230939].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 05 (1080p) [945D7FEB].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 06 (1080p) [E2B6F92C].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 07 (1080p) [2DA7EC8A].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 08 (1080p) [CECA8B90].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 09 (1080p) [1408D11E].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 10 (1080p) [9079C561].mkv")
	input = append(input, "[SubsPlease] Heion Sedai no Idaten-tachi - 11 (1080p) [D298BC5A].mkv")

	expected := make([]EpisodeMetadata, 0, 12)
	expected = append(expected, genOutput("Heion Sedai no Idaten tachi", "%02d", 1, 11)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata12(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, "[Reaktor] The Tatami Galaxy ED [1080p][x265][10-bit].mkv")
	input = append(input, "[Reaktor] The Tatami Galaxy OP [1080p][x265][10-bit].mkv")
	input = append(input, genInput("[Reaktor] The Tatami Galaxy - Special E%d [1080p][x265][10-bit].mkv", 1, 3)...)
	input = append(input, genInput("[Reaktor] The Tatami Galaxy - E%02d [1080p][x265][10-bit].mkv", 1, 11)...)

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genSingle("", "ED"))
	expected = append(expected, genSingle("", "OP"))
	expected = append(expected, genOutput("", "Special E%d", 1, 3)...)
	expected = append(expected, genOutput("", "E%02d", 1, 11)...)

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata13(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, "01. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [E24E9EB2].mkv")
	input = append(input, "02. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [4BF38387].mkv")
	input = append(input, "03. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [1345D78E].mkv")
	input = append(input, "04. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [E189606E].mkv")
	input = append(input, "05. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [236BB482].mkv")
	input = append(input, "06. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [9C258F9E].mkv")
	input = append(input, "07. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [0187A81E].mkv")
	input = append(input, "08. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [5F98DED5].mkv")
	input = append(input, "09. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [DF63FF29].mkv")
	input = append(input, "10. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia] [2222DD11].mkv")
	input = append(input, "11. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia][95080862].mkv")
	input = append(input, "12. Sayonara Zetsubou Sensei [720p Hi10p AAC BDRip][kuchikirukia][799FDE7C].mkv")
	input = append(input, "99. NCED (for linked mkvs and ordered chapters) [06C51833].mkv")
	input = append(input, genInput("Zoku Sayonara Zetsubou Sensei - %02d (BD 1024x576 x264 AAC).mkv", 1, 13)...)
	input = append(input, "Goku Sayonara Zetsubou Sensei - 01 (BD 1280x720 x264 AAC) [EC18711E].mkv")
	input = append(input, "Goku Sayonara Zetsubou Sensei - 02 (BD 1280x720 x264 AAC) [7303476B].mkv")
	input = append(input, "Goku Sayonara Zetsubou Sensei - 03 (BD 1280x720 x264 AAC) [DA5F8783].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 01 (BD 720p h264 AAC) [73898A37].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 02 (BD 720p h264 AAC) [642FBC61].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 03 (BD 720p h264 AAC) [4FA5C3C8].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 04 (BD 720p h264 AAC) [2955BDE3].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 05 (BD 720p h264 AAC) [E5EE016D].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 06 (BD 720p h264 AAC) [161DBF27].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 07 (BD 720p h264 AAC) [997BFEF4].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 08 (BD 720p h264 AAC) [4342F158].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 09 (BD 720p h264 AAC) [C4F54CBF].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 10 (BD 720p h264 AAC) [49BFE205].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 11 (BD 720p h264 AAC) [F7D8331D].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 12 (BD 720p h264 AAC) [1B06902C].mkv")
	input = append(input, "[qcc] Zan Sayonara Zetsubou Sensei - 13 (BD 720p h264 AAC) [F3F90AE9].mkv")
	input = append(input, genInput("Zan Sayonara Zetsubou Sensei Bangaichi - %02d (BD 1280x720 x264 AAC).mkv", 1, 2)...)
	input = append(input, "[Commie] Sayonara Zetsubou Sensei (2012) - BD Special [BD 720p AAC] [BEA51F1F].mkv")

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("Sayonara Zetsubou Sensei", "%02d", 1, 12)...)
	expected = append(expected, genSingle("NCED", "99"))
	expected = append(expected, genOutput("Zoku Sayonara Zetsubou Sensei", "%02d", 1, 13)...)
	expected = append(expected, genOutput("Goku Sayonara Zetsubou Sensei", "%02d", 1, 3)...)
	expected = append(expected, genOutput("Zan Sayonara Zetsubou Sensei", "%02d", 1, 13)...)
	expected = append(expected, genOutput("Zan Sayonara Zetsubou Sensei Bangaichi", "%02d", 1, 2)...)
	expected = append(expected, genSingle("Sayonara Zetsubou Sensei", "BD Special"))

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}

func TestMultipleEpisodeMetadata14(t *testing.T) {
	input := make([]string, 0, 256)
	input = append(input, genInput("[Judas] Hunter x Hunter (2011) - Episodes 001-148/[Judas] Hunter x Hunter (2011) - S01E%03d.mkv", 1, 148)...)
	input = append(input, "[Judas] Hunter x Hunter (2011) - Movies/[Judas] Hunter X Hunter - Movie 1 - Phantom Rouge.mkv")
	input = append(input, "[Judas] Hunter x Hunter (2011) - Movies/[Judas] Hunter X Hunter - Movie 2 - The Last Mission.mkv")

	expected := make([]EpisodeMetadata, 0, 256)
	expected = append(expected, genOutput("01", "%03d", 1, 148)...)
	expected = append(expected, genSingle("Hunter X Hunter", "1"))
	expected = append(expected, genSingle("Hunter X Hunter", "2"))

	metadataArr := ParseMultipleEpisodeMetadata(input)
	assertDiff(t, metadataArr, expected)
}
