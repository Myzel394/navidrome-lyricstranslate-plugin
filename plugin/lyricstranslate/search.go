package lyricstranslate

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/mozillazg/go-unidecode"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

const matchThreshold = 0.85

var (
	searchItemRe = regexp.MustCompile(`(?s)<div class="table__trow">.*?<a href="([^"]+)" class="title-link">\s*<span class="title-text">(.*?)</span></a>.*?<span class="block-1-table__author">\s*<span class="text">(.*?)</span>\s*</span>.*?<div class="table__langs">(.*?)</div>`)
	tagsRe       = regexp.MustCompile(`(?s)<[^>]+>`)
)

func searchForTrack(input lyrics.GetLyricsRequest) (*Song, error) {
	artist := decodeSearchInput(input.Track.Artist)
	title := decodeSearchInput(input.Track.Title)
	normArtist := normalize(artist)
	normTitle := normalize(title)
	query := strings.TrimSpace(normArtist + " " + normTitle)
	endpoint := fmt.Sprintf(utils.LyricstranslateSearchURL, searchPathPart(normArtist), searchPathPart(normTitle))

	utils.LogInfof("searching for '%s' -> %s", query, endpoint)

	body, err := utils.DoGetRequest(endpoint)
	if err != nil || body == nil {
		utils.LogErrorf("search request failed for '%s': %v", query, err)
		return nil, fmt.Errorf("failed to do lyricstranslate search request for query %s; Error: %v", query, err)
	}

	return pickBestMatch(extractSearchHits(string(body)), normArtist, normTitle), nil
}

func extractSearchHits(page string) []Song {
	matches := searchItemRe.FindAllStringSubmatch(page, -1)
	hits := make([]Song, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		hit := Song{
			URL:    absoluteURL(html.UnescapeString(match[1])),
			Title:  cleanHTMLText(match[2]),
			Artist: cleanHTMLText(match[3]),
		}
		if len(match) > 4 {
			hit.Lang = cleanHTMLText(match[4])
		}

		if hit.URL != "" && hit.Title != "" && hit.Artist != "" {
			hits = append(hits, hit)
		}
	}
	return hits
}

func pickBestMatch(hits []Song, normArtist, normTitle string) *Song {
	if bestSong := pickBestMatchWith(hits, normArtist, normTitle, false); bestSong != nil {
		return bestSong
	}
	return pickBestMatchWith(hits, romanize(normArtist), romanize(normTitle), true)
}

func pickBestMatchWith(hits []Song, normArtist, normTitle string, romanized bool) *Song {
	var bestSong *Song
	var bestScore float64

	for _, hit := range hits {
		artistRatio := bestRatio(normArtist, matchVariants(hit.Artist, romanized))
		titleRatio := bestRatio(normTitle, matchVariants(hit.Title, romanized))
		if artistRatio < matchThreshold || titleRatio < matchThreshold {
			continue
		}

		score := (artistRatio + titleRatio) / 2
		if score > bestScore {
			bestScore = score
			bestSong = &hit
		}
	}

	return bestSong
}

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
