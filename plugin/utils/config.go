package utils

import (
	"strconv"

	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
)

func getConfigString(key, def string) string {
	v, ok := pdk.GetConfig(key)
	if !ok || v == "" {
		return def
	}
	return v
}

func ConfigUserAgent() string {
	return getConfigString(ConfigKeyUserAgent, DefaultUserAgent)
}

func ConfigSearchHTTPAcceptHeader() string {
	return getConfigString(ConfigKeyHTTPAccept, DefaultHTTPAccept)
}

func ConfigLevenshteinThreshold() float64 {
	v := getConfigString(ConfigKeyLevenshteinThreshold, strconv.FormatFloat(DefaultLevenshteinThreshold, 'f', -1, 64))
	threshold, err := strconv.ParseFloat(v, 64)
	if err != nil || threshold <= 0 || threshold > 1 {
		return DefaultLevenshteinThreshold
	}
	return threshold
}
