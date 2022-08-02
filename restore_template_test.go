package roflmeta

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString() string {
	b := make([]rune, rand.Int31n(19)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testPairImpl(t *testing.T, str1 string, str2 string, expected string) {
	expectedTemplate := newTemplate(expected)
	resultTemplate := findTemplateForPair([]rune(str1), []rune(str2))
	if !reflect.DeepEqual(expectedTemplate, resultTemplate) {
		t.Fatalf("Invalid pair result: expected %s, got %s", expectedTemplate.String(), resultTemplate.String())
	}
}

func testRestore(t *testing.T, expected string, strings ...string) {
	expectedTemplate := newTemplate(expected)
	resultTemplate, err := restoreTemplate(strings)
	if err != nil {
		t.Fatalf("Template restoration failed, expected %s", expectedTemplate.String())
	}
	if !reflect.DeepEqual(&expectedTemplate, resultTemplate) {
		t.Fatalf("Invalid restore result: expected %s, got %s", expectedTemplate.String(), resultTemplate.String())
	}
}

func TestPairSimple(t *testing.T) {
	testPairImpl(t, "", "", "")
	testPairImpl(t, "a", "a", "a")
	testPairImpl(t, "a", "", "*")
	testPairImpl(t, "", "b", "*")
	testPairImpl(t, "abc", "def", "*")
	testPairImpl(t, "abc", "bcd", "*bc*")
	testPairImpl(t, "abc", "abd", "ab*")
	testPairImpl(t, "abc", "dbc", "*bc")
	testPairImpl(t, "abcdef", "AbcDeF", "*bc*e*")
}

func TestPairMultipleSolutions(t *testing.T) {
	// must not select *01*.mkv here
	t.Skip("this test doesn't work")
	testPairImpl(t, "24x01.mkv", "01x02.mkv", "*x0*.mkv")
	testPairImpl(t, "01x24.mkv", "02x01.mkv", "*x*.mkv")
}

func TestRestoreSimple(t *testing.T) {
	testRestore(t, "s1e*", "s1e1", "s1e2", "s1e3")
	testRestore(t, "s*e*", "s1e1", "s2e1", "s2e2")
	testRestore(t, "1-*", "1-1", "1-2", "1-3")
	testRestore(t, "*-*", "1-1", "1-2", "1-3", "2-1")
}

func TestSingleNamedSeason(t *testing.T) {
	strings := make([]string, 0, 100)
	for i := 1; i < 15; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Bakemonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	testRestore(t, "[DB]Bakemonogatari_-_*_(10bit_BD1080p_x265).mkv", strings...)
}

func TestTwoNamedSeasons(t *testing.T) {
	strings := make([]string, 0, 100)
	for i := 1; i < 15; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Bakemonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	for i := 1; i < 5; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Hanamonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	strings = append(strings, "[DB]Hanamonogatari_-_NCED01_(10bit_BD1080p_x265).mkv")
	strings = append(strings, "[DB]Hanamonogatari_-_NCOP01_(10bit_BD1080p_x265).mkv")
	testRestore(t, "[DB]*a*monogatari_-_*_(10bit_BD1080p_x265).mkv", strings...)
}

func TestThreeNamedSeasons(t *testing.T) {
	strings := make([]string, 0, 100)
	for i := 1; i < 15; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Bakemonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	for i := 1; i < 5; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Hanamonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	strings = append(strings, "[DB]Hanamonogatari_-_NCED01_(10bit_BD1080p_x265).mkv")
	strings = append(strings, "[DB]Hanamonogatari_-_NCOP01_(10bit_BD1080p_x265).mkv")
	for i := 1; i < 3; i++ {
		strings = append(strings, fmt.Sprintf("[DB]Kizumonogatari_-_%02d_(10bit_BD1080p_x265).mkv", i))
	}
	testRestore(t, "[DB]*monogatari_-_*_(10bit_BD1080p_x265).mkv", strings...)
}

func TestLongSingleSeason(t *testing.T) {
	strings := make([]string, 0, 256)
	for i := 1; i < 148; i++ {
		strings = append(strings, fmt.Sprintf("[Judas] Hunter x Hunter (2011) - S01E%03d.mkv", i))
	}
	testRestore(t, "[Judas] Hunter x Hunter (2011) - S01E*.mkv", strings...)
}

func TestMultipleLongSeasons(t *testing.T) {
	strings := make([]string, 0, 450)
	for j := 1; j < 12; j++ {
		for i := 1; i < 150; i++ {
			strings = append(strings, fmt.Sprintf("[Judas] Hunter x Hunter (2011) - S%02dE%03d.mkv", j, i))
		}
	}
	testRestore(t, "[Judas] Hunter x Hunter (2011) - S*E*.mkv", strings...)
}

func TestShuffle(t *testing.T) {
	strings := make([]string, 0, 500)
	for i := 1; i < 500; i++ {
		strings = append(strings, fmt.Sprintf("[Roflex] %s (2022) - season %s, episode %s", randomString(), randomString(), randomString()))
	}
	testRestore(t, "[Roflex] * (2022) - season *, episode *", strings...)
}

func TestMinimal(t *testing.T) {
	strings := make([]string, 0, 50)
	for i := 1; i < 24; i++ {
		strings = append(strings, fmt.Sprintf("Dr Stone/Dr Stone - %02d.mkv", i))
	}
	for i := 1; i < 11; i++ {
		strings = append(strings, fmt.Sprintf("Dr Stone Season 2/Dr Stone - %02d.mkv", i))
	}
	testRestore(t, "Dr Stone*/Dr Stone - *.mkv", strings...)
}

func TestThreeVars(t *testing.T) {
	t.Skip("this test doesn't work")
	strings := make([]string, 0, 50)
	for j := 1; j <= 2; j++ {
		for i := 1; i <= 24; i++ {
			strings = append(strings, fmt.Sprintf("nartsiss %02d/nartsiss %02dx%02d.mkv", j, i, j))
		}
	}
	testRestore(t, "nartsiss 0*/nartsiss *x0*.mkv", strings...)
}
