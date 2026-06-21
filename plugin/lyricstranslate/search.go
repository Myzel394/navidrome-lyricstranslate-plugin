package lyricstranslate

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

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

	return pickBestMatch(extractSearchHits(string(body)), normArtist, normTitle, utils.ConfigLevenshteinThreshold()), nil
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

func pickBestMatch(hits []Song, normArtist, normTitle string, matchThreshold float64) *Song {
	if bestSong := pickBestMatchWith(hits, normArtist, normTitle, false, matchThreshold); bestSong != nil {
		return bestSong
	}
	return pickBestMatchWith(hits, romanize(normArtist), romanize(normTitle), true, matchThreshold)
}

func pickBestMatchWith(hits []Song, normArtist, normTitle string, romanized bool, matchThreshold float64) *Song {
	var bestSong *Song
	var bestScore float64

	for _, hit := range hits {
		artistRatio := bestRatio(normArtist, artistMatchVariants(hit.Artist, romanized))
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
