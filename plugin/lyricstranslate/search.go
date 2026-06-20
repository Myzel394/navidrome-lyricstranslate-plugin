package lyricstranslate

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

const matchThreshold = 0.85

var (
	songsBlockRe = regexp.MustCompile(`(?s)<div class="song-list search-res__block block-search-res _songs">(.*?)<a href="/en/songs/`)
	searchItemRe = regexp.MustCompile(`(?s)<div class="block-search-res__item">.*?<a href="([^"]+)" class="block-1-table__title"><h2>(.*?)</h2></a>.*?<a href="[^"]+" class="block-1-table__author"><span class="text">(.*?)</span></a>.*?<div class="table__langs">(.*?)</div>`)
	tagsRe       = regexp.MustCompile(`(?s)<[^>]+>`)
)

func searchForTrack(input lyrics.GetLyricsRequest) (*Song, error) {
	normArtist := normalize(input.Track.Artist)
	normTitle := normalize(input.Track.Title)
	query := strings.TrimSpace(normArtist + " " + normTitle)
	endpoint := fmt.Sprintf(utils.LyricstranslateSearchURL, url.QueryEscape(query))

	utils.LogInfof("searching for '%s' -> %s", query, endpoint)

	body, err := utils.DoGetRequest(endpoint)
	if err != nil || body == nil {
		utils.LogErrorf("search request failed for '%s': %v", query, err)
		return nil, fmt.Errorf("failed to do lyricstranslate search request for query %s; Error: %v", query, err)
	}

	return pickBestMatch(extractSearchHits(string(body)), normArtist, normTitle), nil
}

func extractSearchHits(page string) []Song {
	blockMatch := songsBlockRe.FindStringSubmatch(page)
	if len(blockMatch) < 2 {
		return nil
	}

	matches := searchItemRe.FindAllStringSubmatch(blockMatch[1], -1)
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
	var bestSong *Song
	var bestScore float64

	for _, hit := range hits {
		hitArtist := normalize(hit.Artist)
		hitTitle := normalize(hit.Title)

		artistRatio := levenshteinRatio(normArtist, hitArtist)
		titleRatio := levenshteinRatio(normTitle, hitTitle)
		if artistRatio < matchThreshold || titleRatio < matchThreshold {
			continue
		}

		score := (artistRatio + titleRatio) / 2
		if score > bestScore {
			bestScore = score
			hit.Artist = hitArtist
			hit.Title = hitTitle
			bestSong = &hit
		}
	}

	return bestSong
}

func cleanHTMLText(s string) string {
	s = tagsRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	s = whitespaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
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
