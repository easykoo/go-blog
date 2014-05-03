package common

import (
	"fmt"
	cfg "github.com/Unknwon/goconfig"
	"strings"
)

var Tsl *cfg.ConfigFile

func init() {
	Tsl, _ = cfg.LoadConfigFile("messages.ini")
}

func Translate(lang string, format string) string {
	if lang == "" || !strings.Contains(lang, "zh") {
		lang = "en"
	}
	return Tsl.MustValue(lang, format, format)
}

func Translatef(lang string, format string, args ...interface{}) string {
	if lang == "" || !strings.Contains(lang, "zh") {
		lang = "en"
	}
	return fmt.Sprintf(Tsl.MustValue(lang, format, format), args)
}
