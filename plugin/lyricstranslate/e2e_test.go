package lyricstranslate

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Myzel394/navidrome-lyricstranslate-plugin/plugin/utils"
	"github.com/navidrome/navidrome/plugins/pdk/go/lyrics"
)

func TestFetchLyricsNeverGonnaGiveYouUpEndToEnd(t *testing.T) {
	const (
		artist = "Rick Astley"
		title  = "Never Gonna Give You Up"
	)
	mockHTTP(t)

	resp, err := FetchLyrics(lyrics.GetLyricsRequest{
		Track: lyrics.TrackInfo{
			Artist: artist,
			Title:  title,
		},
	})
	if err != nil {
		t.Fatalf("FetchLyrics returned error: %v", err)
	}
	if len(resp.Lyrics) == 0 {
		t.Fatal("expected at least one lyrics result")
	}

	got := resp.Lyrics[0]
	if got.Lang != "en" {
		t.Fatalf("expected language en, got %q", got.Lang)
	}
	if !strings.Contains(got.Text, "Never gonna give you up") {
		t.Fatalf("expected fetched lyrics to contain chorus, got %q", got.Text)
	}
	if !strings.Contains(got.Text, "[00:") {
		t.Fatalf("expected synced LRC lyrics, got %q", got.Text)
	}
}

func mockHTTP(t *testing.T) {
	t.Helper()

	originalDoGetRequest := utils.DoGetRequest
	originalLogInfof := utils.LogInfof
	originalLogErrorf := utils.LogErrorf
	t.Cleanup(func() {
		utils.DoGetRequest = originalDoGetRequest
		utils.LogInfof = originalLogInfof
		utils.LogErrorf = originalLogErrorf
	})

	utils.DoGetRequest = realGetRequest
	utils.LogInfof = t.Logf
	utils.LogErrorf = t.Logf
}

func realGetRequest(endpoint string) ([]byte, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", utils.DefaultHTTPAccept)
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("User-Agent", utils.DefaultUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != utils.HTTPStatusOK {
		return body, fmt.Errorf("error code %d returned from Lyricstranslate for endpoint %s", resp.StatusCode, endpoint)
	}
	return body, nil
}
