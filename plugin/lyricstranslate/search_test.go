package lyricstranslate

import "testing"

func TestPickBestMatchUsesRomanizedFallback(t *testing.T) {
	hits := []Song{{
		Artist: "Tamara Grujeska",
		Title:  "Кажи ми, кажи ми кој",
		URL:    "https://lyricstranslate.com/en/tamara-grujeska-kazhi-mi-kazhi-mi-koj-lyrics.html",
	}}

	match := pickBestMatch(hits, normalize("Tamara Grujeska"), normalize("Kazi Mi, Kazi Mi Koj"))
	if match == nil {
		t.Fatal("expected romanized fallback match")
	}
	if match.URL != hits[0].URL {
		t.Fatalf("unexpected match URL: %s", match.URL)
	}
}

func TestPickBestMatchUsesBracketedRomanizedVariant(t *testing.T) {
	hits := []Song{{
		Artist: "Tamara Grujeska",
		Title:  "Кажи ми, кажи ми кој (Kaži mi, kaži mi koj)",
		URL:    "https://lyricstranslate.com/en/tamara-grujeska-kazhi-mi-kazhi-mi-koj-lyrics.html",
	}}

	match := pickBestMatch(hits, normalize("Tamara Grujeska"), normalize("Kazi Mi, Kazi Mi Koj"))
	if match == nil {
		t.Fatal("expected bracketed romanized variant match")
	}
	if match.URL != hits[0].URL {
		t.Fatalf("unexpected match URL: %s", match.URL)
	}
}

func TestDecodeSearchInputHandlesDoubleEncoding(t *testing.T) {
	got := decodeSearchInput("Sep%2520%2526%2520Jasmijn")
	want := "Sep & Jasmijn"
	if got != want {
		t.Fatalf("decodeSearchInput() = %q, want %q", got, want)
	}
}
