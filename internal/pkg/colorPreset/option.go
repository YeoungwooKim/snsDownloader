package colorPreset

import (
	"runtime"
	"strings"
)

//text options
var (
	Reset, BoldOn, BoldOff, UnderLineOn, UnderLineOff string
)

// text colors
var (
	LightGray, LightRed, LightGreen, LightYellow, LightBlue, LightMagenta, LightCyan, LightWhite string
	Black, Red, Green, Yellow, Blue, Magenta, Cyan, White                                        string
	Default                                                                                      string
)

// background colors
var (
	BgLightGray, BgLightRed, BgLightGreen, BgLightYellow, BgLightBlue, BgLightMagenta, BgLightCyan, BgLightWhite string
	BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite                                        string
	BgDefault                                                                                                    string
)

func init() {
	osName := runtime.GOOS
	if !strings.Contains(strings.ToLower(osName), "window") {
		Reset = "\x1b[0m"
		BoldOn = "\x1b[1m"
		BoldOff = "\x1b[21m"
		UnderLineOn = "\x1b[4m"
		UnderLineOff = "\x1b[24m"

		Black = "\x1b[30m"
		Red = "\x1b[31m"
		Green = "\x1b[32m"
		Yellow = "\x1b[33m"
		Blue = "\x1b[34m"
		Magenta = "\x1b[35m"
		Cyan = "\x1b[36m"
		White = "\x1b[37m"
		Default = "\x1b[39m"
		LightGray = "\x1b[90m"
		LightRed = "\x1b[91m"
		LightGreen = "\x1b[92m"
		LightYellow = "\x1b[93m"
		LightBlue = "\x1b[94m"
		LightMagenta = "\x1b[95m"
		LightCyan = "\x1b[96m"
		LightWhite = "\x1b[97m"

		BgBlack = "\x1b[40m"
		BgRed = "\x1b[41m"
		BgGreen = "\x1b[42m"
		BgYellow = "\x1b[43m"
		BgBlue = "\x1b[44m"
		BgMagenta = "\x1b[45m"
		BgCyan = "\x1b[46m"
		BgWhite = "\x1b[47m"
		BgDefault = "\x1b[49m"
		BgLightGray = "\x1b[100m"
		BgLightRed = "\x1b[101m"
		BgLightGreen = "\x1b[102m"
		BgLightYellow = "\x1b[103m"
		BgLightBlue = "\x1b[104m"
		BgLightMagenta = "\x1b[105m"
		BgLightCyan = "\x1b[106m"
		BgLightWhite = "\x1b[107m"
	}
}
