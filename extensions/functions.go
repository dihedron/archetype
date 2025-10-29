package extensions

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"include":  Include,
		"dump":     DumpArgs,
		"api":      CallAPI,
		"fileSize": FileSize,
		"dirSize":  DirSize,
		"isFile":   IsFile,
		"isDir":    IsDir,
		"listDir":  ListDir,
		// "fetch":     FetchHTML,
		// "html":      HTML,
		"blue":      Blue,
		"cyan":      Cyan,
		"green":     Green,
		"magenta":   Magenta,
		"purple":    Magenta,
		"red":       Red,
		"yellow":    Yellow,
		"white":     White,
		"hiblue":    HighBlue,
		"hicyan":    HighCyan,
		"higreen":   HighGreen,
		"himagenta": HighMagenta,
		"hipurple":  HighMagenta,
		"hired":     HighRed,
		"hiyellow":  HighYellow,
		"hiwhite":   HighWhite,
	}
}
