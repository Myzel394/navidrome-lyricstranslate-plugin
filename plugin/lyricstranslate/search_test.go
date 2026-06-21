package lyricstranslate

import "testing"

func TestPickBestMatchHandlesStageNamePrefix(t *testing.T) {
	hits := []Song{{
		Artist: "Arhanna Sandra Arbma",
		Title:  "Hoiame kokku",
		URL:    "https://lyricstranslate.com/en/arhanna-sandra-arbma-hoiame-kokku-lyrics.html",
	}}

	match := pickBestMatch(hits, normalize("ARHANNA"), normalize("Hoiame Kokku"), 0.75)
	if match == nil {
		t.Fatal("expected stage-name artist match")
	}
	if match.URL != hits[0].URL {
		t.Fatalf("unexpected match URL: %s", match.URL)
	}
}
