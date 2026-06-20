package main

import (
	"fmt"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/lyricstranslate"
	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

func (p *plugin) GetLyrics(input lyrics.GetLyricsRequest) (lyrics.GetLyricsResponse, error) {
	resp, err := lyricstranslate.FetchLyrics(input)
	if err != nil {
		utils.LogErrorf("GetLyrics failed: %v", err)
		return resp, fmt.Errorf("%s%w", utils.LogPrefix, err)
	}
	return resp, nil
}
