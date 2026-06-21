package lyricstranslate

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	bracketsRe       = regexp.MustCompile(`[\(\[\{][^\)\]\}]*[\)\]\}]`)
	bracketContentRe = regexp.MustCompile(`[\(\[\{]([^\)\]\}]*)[\)\]\}]`)
	dashSuffixRe     = regexp.MustCompile(`(?i)\s*-\s*(remaster(ed)?|single version|live|deluxe|edit|mix|version|radio edit|extended).*$`)
	artistJoinersRe  = regexp.MustCompile(`(?i)\s*(?:&|/|,|\band\b|\bfeat\.?\b|\bft\.?\b|\bwith\b|\bx\b)\s*`)
	whitespaceRe     = regexp.MustCompile(`\s+`)
)

func normalize(s string) string {
	s = strings.ToLower(s)
	s = stripDiacritics(s)
	s = bracketsRe.ReplaceAllString(s, " ")
	s = dashSuffixRe.ReplaceAllString(s, "")
	s = whitespaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func stripDiacritics(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}

func levenshteinRatio(s1, s2 string) float64 {
	r1, r2 := []rune(s1), []rune(s2)
	len1, len2 := len(r1), len(r2)

	if len1 == 0 && len2 == 0 {
		return 1.0
	}
	if len1 == 0 || len2 == 0 {
		return 0.0
	}

	d := make([][]int, len1+1)
	for i := range d {
		d[i] = make([]int, len2+1)
		d[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		d[0][j] = j
	}

	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}

			min := d[i-1][j] + 1
			if d[i][j-1]+1 < min {
				min = d[i][j-1] + 1
			}
			if d[i-1][j-1]+cost < min {
				min = d[i-1][j-1] + cost
			}
			d[i][j] = min
		}
	}

	maxLen := float64(len1)
	if float64(len2) > maxLen {
		maxLen = float64(len2)
	}

	return 1.0 - (float64(d[len1][len2]) / maxLen)
}
