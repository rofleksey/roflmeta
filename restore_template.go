package roflmeta

import (
	"regexp"
	"strings"
)

type template struct {
	runes []rune
}

func newTemplate(s string) template {
	return template{[]rune(s)}
}

func (t *template) varCount() int {
	result := 0
	for _, r := range t.runes {
		if r == '*' {
			result++
		}
	}
	return result
}

func (t *template) check(filenames []string) error {
	regex := t.toRegex()
	for _, s := range filenames {
		test := regex.FindStringSubmatch(s)
		if test == nil {
			return errInvalidTemplate
		}
	}
	return nil
}

// fix searches and removes redundant vars
func (t *template) fix(filenames []string) template {
	regex := t.toRegex()
	varCount := t.varCount()
	hasNonEmpty := make([]int, varCount)
	for _, s := range filenames {
		test := regex.FindStringSubmatch(s)
		if test == nil {
			return *t
		}
		for i := 1; i <= varCount; i++ {
			if !spaceRegex.MatchString(test[i]) {
				hasNonEmpty[i-1]++
			}
		}
	}
	varsToRemove := make([]int, 0, varCount)
	for i, value := range hasNonEmpty {
		if value == 0 {
			varsToRemove = append(varsToRemove, i)
		}
	}
	if len(varsToRemove) > 0 {
		return t.removeVars(varsToRemove)
	}
	return *t
}

func (t *template) removeVars(indices []int) template {
	if len(indices) == 0 {
		return *t
	}
	runes := make([]rune, 0, len(t.runes)-len(indices))
	curIndex := -1
	removeIndex := 0
	for _, r := range t.runes {
		if r == '*' {
			curIndex++
			if removeIndex < len(indices) && curIndex == indices[removeIndex] {
				removeIndex++
				continue
			}
		}
		runes = append(runes, r)
	}
	return template{runes}
}

func (t *template) toRegexString() string {
	var builder strings.Builder
	var subBuilder strings.Builder
	builder.WriteRune('^')
	for _, r := range t.runes {
		if r == '*' {
			builder.WriteString(regexp.QuoteMeta(subBuilder.String()))
			subBuilder.Reset()
			builder.WriteString("(.*?)")
		} else {
			subBuilder.WriteRune(r)
		}
	}
	if subBuilder.Len() > 0 {
		builder.WriteString(regexp.QuoteMeta(subBuilder.String()))
	}
	builder.WriteRune('$')
	return builder.String()
}

func (t *template) toRegex() *regexp.Regexp {
	regexString := t.toRegexString()
	return regexp.MustCompile(regexString)
}

func (t *template) merge(other template) template {
	result := findTemplateForPair(t.runes, other.runes)
	runes := make([]rune, 0, len(result.runes))
	lastWasVar := false
	for _, r := range result.runes {
		if r == '*' {
			if lastWasVar {
				continue
			}
			lastWasVar = true
		} else {
			lastWasVar = false
		}
		runes = append(runes, r)
	}
	return template{runes}
}

func (t *template) String() string {
	return string(t.runes)
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// findTemplateForPair is actually the Longest Common Subsequence algorithm with stars insertion on each skip
// also it reverses the strings to find the last possible LCS
func findTemplateForPair(str1 []rune, str2 []rune) template {
	len1 := len(str1)
	len2 := len(str2)
	table := make([][]int, len1+1)
	for i := range table {
		table[i] = make([]int, len2+1)
	}
	for i := 0; i <= len1; i++ {
		for j := 0; j <= len2; j++ {
			if i == 0 || j == 0 {
				table[i][j] = 0
			} else if str1[i-1] == str2[j-1] {
				table[i][j] = table[i-1][j-1] + 1
			} else {
				table[i][j] = max(table[i-1][j], table[i][j-1])
			}
		}
	}
	index := table[len1][len2]
	resultRunes := make([]rune, index)
	skips := make([]bool, index+1)
	skipsCount := 0
	i1 := len1
	j1 := len2
	for i1 > 0 && j1 > 0 {
		if str1[i1-1] == str2[j1-1] {
			resultRunes[index-1] = str1[i1-1]
			i1--
			j1--
			index--
		} else if table[i1-1][j1] > table[i1][j1-1] {
			i1--
			if !skips[index] {
				skips[index] = true
				skipsCount++
			}
		} else {
			j1--
			if !skips[index] {
				skips[index] = true
				skipsCount++
			}
		}
	}
	if i1 > 0 || j1 > 0 {
		if !skips[0] {
			skips[0] = true
			skipsCount++
		}
	}
	result := make([]rune, len(resultRunes)+skipsCount)
	resultI := 0
	if len(resultRunes) == 0 {
		if skips[0] {
			result[0] = '*'
		}
	} else {
		for i := 0; i < len(resultRunes); i++ {
			if skips[i] {
				result[resultI] = '*'
				resultI++
			}
			result[resultI] = resultRunes[i]
			resultI++
		}
		if skips[len(resultRunes)] {
			result[resultI] = '*'
		}
	}
	return template{result}
}

func restoreTemplate(filenames []string) (*template, error) {
	if len(filenames) == 0 {
		return &template{}, nil
	}
	if len(filenames) == 1 {
		return &template{[]rune(filenames[0])}, nil
	}
	curTemplate := findTemplateForPair([]rune(filenames[0]), []rune(filenames[1]))
	for i := 2; i < len(filenames); i++ {
		pairTemplate := findTemplateForPair([]rune(filenames[i-1]), []rune(filenames[i]))
		curTemplate = curTemplate.merge(pairTemplate)
	}
	err := curTemplate.check(filenames)
	if err != nil {
		return nil, err
	}
	curTemplate = curTemplate.fix(filenames)
	return &curTemplate, nil
}
