package utils

const LogPrefix = "navidrome-lyricstranslate-plugin: "

const HTTPStatusOK = 200

const (
	DefaultUserAgent  = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Mobile/15E148 Safari/604.1"
	DefaultHTTPAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
)

const (
	ConfigKeyUserAgent  = "lyricstranslate_user_agent"
	ConfigKeyHTTPAccept = "lyricstranslate_http_accept"
)

const (
	LyricstranslateBaseURL      = "https://lyricstranslate.com"
	LyricstranslateSearchURL    = LyricstranslateBaseURL + "/en/songs/0/%s/%s/0/none/0?order=relevance"
	LyricstranslateSubtitlesURL = LyricstranslateBaseURL + "/en/callback/subtitles/%s/%s"
)
