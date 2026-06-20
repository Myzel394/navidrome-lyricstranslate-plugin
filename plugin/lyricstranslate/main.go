package lyricstranslate

import (
	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

type Song struct {
	Artist string
	Title  string
	URL    string
}

const stubLyrics = "These are stub lyrics from the Lyricstranslate plugin."

func FetchLyrics(input lyrics.GetLyricsRequest) (lyrics.GetLyricsResponse, error) {
	utils.LogInfof("FetchLyrics stub: artist='%s' title='%s'", input.Track.Artist, input.Track.Title)
	return lyrics.GetLyricsResponse{
		Lyrics: []lyrics.LyricsText{{Lang: "en", Text: stubLyrics}},
	}, nil
}
