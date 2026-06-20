package utils

import (
	"fmt"

	"github.com/navidrome/navidrome/plugins/pdk/go/pdk"
)

var LogInfof = logInfof

var LogErrorf = logErrorf

func logInfof(format string, args ...any) {
	pdk.Log(pdk.LogInfo, fmt.Sprintf(LogPrefix+format, args...))
}

func logErrorf(format string, args ...any) {
	pdk.Log(pdk.LogError, fmt.Sprintf(LogPrefix+format, args...))
}
