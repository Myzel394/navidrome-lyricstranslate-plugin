package lyricstranslate

import (
	"html"
	"net/url"
	"strings"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/mozillazg/go-unidecode"
)

func bestRatio(input string, variants []string) float64 {
	var best float64
	for _, variant := range variants {
		ratio := levenshteinRatio(input, variant)
		if ratio > best {
			best = ratio
		}
	}
	return best
}

func matchVariants(s string, romanized bool) []string {
	variants := []string{s, bracketsRe.ReplaceAllString(s, " ")}
	for _, match := range bracketContentRe.FindAllStringSubmatch(s, -1) {
		if len(match) > 1 {
			variants = append(variants, match[1])
		}
	}
	return normalizeVariants(variants, romanized)
}

func artistMatchVariants(s string, romanized bool) []string {
	variants := []string{s, bracketsRe.ReplaceAllString(s, " ")}
	for _, part := range artistJoinersRe.Split(s, -1) {
		variants = append(variants, part)
	}
	for _, part := range strings.Fields(s) {
		if len([]rune(part)) > 3 {
			variants = append(variants, part)
		}
	}
	return normalizeVariants(variants, romanized)
}

func normalizeVariants(variants []string, romanized bool) []string {
	seen := make(map[string]struct{}, len(variants))
	normalized := make([]string, 0, len(variants))
	for _, variant := range variants {
		if romanized {
			variant = romanize(variant)
		}
		variant = normalize(variant)
		if variant == "" {
			continue
		}
		if _, ok := seen[variant]; ok {
			continue
		}
		seen[variant] = struct{}{}
		normalized = append(normalized, variant)
	}
	return normalized
}

func romanize(s string) string {
	return normalize(unidecode.Unidecode(s))
}

func cleanHTMLText(s string) string {
	s = tagsRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	s = whitespaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func decodeSearchInput(s string) string {
	s = html.UnescapeString(strings.TrimSpace(s))
	for i := 0; i < 2; i++ {
		decoded, err := url.QueryUnescape(s)
		if err != nil || decoded == s {
			break
		}
		s = decoded
	}
	return s
}

func searchPathPart(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "none"
	}
	return strings.ReplaceAll(url.PathEscape(s), "%26", "&")
}

func absoluteURL(rawURL string) string {
	if rawURL == "" || strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	if strings.HasPrefix(rawURL, "/") {
		return utils.LyricstranslateBaseURL + rawURL
	}
	return utils.LyricstranslateBaseURL + "/" + rawURL
}
