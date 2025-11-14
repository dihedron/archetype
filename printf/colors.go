package printf

import "github.com/fatih/color"

var (
	Green   func(...any) string = color.New(color.FgGreen).SprintFunc()
	Red     func(...any) string = color.New(color.FgRed).SprintFunc()
	Yellow  func(...any) string = color.New(color.FgYellow).SprintFunc()
	Blue    func(...any) string = color.New(color.FgBlue).SprintFunc()
	Magenta func(...any) string = color.New(color.FgMagenta).SprintFunc()
	White   func(...any) string = color.New(color.FgWhite).SprintFunc()
)
