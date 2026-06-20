package lyricstranslate

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	astisub "github.com/asticode/go-astisub"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

type subtitlesResponse struct {
	WebVTT   string `json:"webtvv"`
	LangCode string `json:"langsrc"`
	LangName string `json:"langname"`
}

var (
	songBodyRe       = regexp.MustCompile(`(?s)<div class="translate__text[^"]*" id="song-body">(.*?)<div id="song-transliteration">`)
	subtitlesAttrsRe = regexp.MustCompile(`(?s)<span class=['"]video-icon-player[^'"]*['"][^>]*\bnid=['"]([^'"]+)['"][^>]*\byoutube=['"]([^'"]+)['"][^>]*\bcc=['"]1['"]`)
	adBlockRe        = regexp.MustCompile(`(?s)<div id="adngin-[^>]*>.*?</div>`)
	lineBreakTagRe   = regexp.MustCompile(`(?i)<br\s*/?>`)
	emptyLineRe      = regexp.MustCompile(`(?s)<div class=['"]emptyline['"][^>]*>.*?</div>`)
	lyricLineCloseRe = regexp.MustCompile(`(?s)</div>\s*<div class="ll-[^"]*">`)
	stripTagsRe      = regexp.MustCompile(`(?s)<[^>]+>`)
	numberLineRe     = regexp.MustCompile(`^\d+$`)
	vttTimingLineRe  = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}[\.,]\d{3}\s+-->\s+\d{2}:\d{2}:\d{2}[\.,]\d{3}`)
	blankLinesRe     = regexp.MustCompile(`\n{3,}`)
)

func fetchLyricsForTrack(track *Song) (lyrics.GetLyricsResponse, error) {
	body, err := utils.DoGetRequest(track.URL)
	if err != nil || body == nil {
		return lyrics.GetLyricsResponse{}, fmt.Errorf("failed to fetch lyricstranslate page for %s; Error: %v", track.URL, err)
	}
	page := string(body)

	if text, lang, synced, err := fetchLyricsFromSubtitles(page); err == nil && text != "" {
		if !synced {
			utils.LogInfof("found plain lyrics from subtitles for %s", track.URL)
			return lyrics.GetLyricsResponse{
				Lyrics: []lyrics.LyricsText{{Lang: lang, Text: text}},
			}, nil
		}

		utils.LogInfof("found synced lyrics for %s", track.URL)
		return lyrics.GetLyricsResponse{
			Lyrics: []lyrics.LyricsText{{Lang: lang, Text: text}},
		}, nil
	}

	text, err := extractLyricsFromHTML(page)
	if err != nil {
		return lyrics.GetLyricsResponse{}, err
	}

	utils.LogInfof("found lyrics for %s", track.URL)
	return lyrics.GetLyricsResponse{
		Lyrics: []lyrics.LyricsText{{Lang: languageCode(track.Lang), Text: text}},
	}, nil
}

func fetchLyricsFromSubtitles(page string) (string, string, bool, error) {
	match := subtitlesAttrsRe.FindStringSubmatch(page)
	if len(match) < 3 {
		return "", "", false, fmt.Errorf("could not find subtitle video attributes")
	}

	endpoint := fmt.Sprintf(utils.LyricstranslateSubtitlesURL, match[1], match[2])
	body, err := utils.DoGetRequest(endpoint)
	if err != nil || body == nil {
		return "", "", false, fmt.Errorf("failed to fetch lyricstranslate subtitles for %s/%s; Error: %v", match[1], match[2], err)
	}

	var resp subtitlesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", "", false, fmt.Errorf("failed to parse lyricstranslate subtitles response: %v", err)
	}
	if resp.WebVTT == "" {
		return "", "", false, fmt.Errorf("empty lyricstranslate subtitles response")
	}

	lang := resp.LangCode
	if lang == "" {
		lang = languageCode(resp.LangName)
	}

	lrc, plain := parseWebVTT(resp.WebVTT)
	if lrc != "" {
		return lrc, lang, true, nil
	}
	if plain != "" {
		return plain, lang, false, nil
	}

	return "", "", false, fmt.Errorf("no lyricstranslate subtitle text found")
}

func parseWebVTT(vtt string) (string, string) {
	subs, err := astisub.ReadFromWebVTT(strings.NewReader(vtt))
	if err != nil || subs == nil || len(subs.Items) == 0 {
		return "", plainTextFromRawWebVTT(vtt)
	}

	lrcLines := make([]string, 0, len(subs.Items))
	plainLines := make([]string, 0, len(subs.Items))

	for _, item := range subs.Items {
		if item == nil {
			continue
		}

		text := strings.TrimSpace(item.String())
		if text == "" {
			continue
		}

		plainLines = append(plainLines, text)
		lrcLines = append(lrcLines, durationToLRC(item.StartAt)+singleLineSubtitleText(text))
	}

	return strings.Join(lrcLines, "\n"), strings.Join(plainLines, "\n")
}

func durationToLRC(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	totalHundredths := int(d / (10 * time.Millisecond))
	minutes := totalHundredths / 6000
	seconds := (totalHundredths / 100) % 60
	hundredths := totalHundredths % 100
	return fmt.Sprintf("[%02d:%02d.%02d]", minutes, seconds, hundredths)
}

func singleLineSubtitleText(text string) string {
	lines := strings.Split(text, "\n")
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			parts = append(parts, line)
		}
	}
	return strings.Join(parts, " / ")
}

func plainTextFromRawWebVTT(vtt string) string {
	lines := strings.Split(vtt, "\n")
	plain := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "WEBVTT" || numberLineRe.MatchString(line) || vttTimingLineRe.MatchString(line) {
			continue
		}

		line = stripTagsRe.ReplaceAllString(line, "")
		line = html.UnescapeString(line)
		line = strings.TrimSpace(line)
		if line != "" {
			plain = append(plain, line)
		}
	}

	return strings.Join(plain, "\n")
}

func extractLyricsFromHTML(page string) (string, error) {
	match := songBodyRe.FindStringSubmatch(page)
	if len(match) < 2 {
		return "", fmt.Errorf("could not find song-body lyrics block")
	}

	text := match[1]
	text = adBlockRe.ReplaceAllString(text, "")
	text = lineBreakTagRe.ReplaceAllString(text, "\n")
	text = emptyLineRe.ReplaceAllString(text, "\n\n")
	text = lyricLineCloseRe.ReplaceAllString(text, "\n")
	text = stripTagsRe.ReplaceAllString(text, "")
	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = normalizeLines(text)

	if text == "" {
		return "", fmt.Errorf("lyrics text not found in song-body block")
	}
	return text, nil
}

func normalizeLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	text = strings.Join(lines, "\n")
	text = blankLinesRe.ReplaceAllString(text, "\n\n")
	return strings.TrimSpace(text)
}

func languageCode(lang string) string {
	switch strings.ToLower(strings.TrimSpace(lang)) {
	case "english":
		return "en"
	case "german":
		return "de"
	case "spanish":
		return "es"
	case "french":
		return "fr"
	case "italian":
		return "it"
	case "japanese":
		return "ja"
	case "portuguese":
		return "pt"
	case "russian":
		return "ru"
	default:
		return ""
	}
}
