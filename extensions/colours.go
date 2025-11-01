package extensions

import (
	"github.com/jedib0t/go-pretty/v6/text"
)

// HighBlue returns the string representation of the given value in high-intensity blue.
func HighBlue(v interface{}) string {
	return text.FgHiBlue.Sprint(v)
}

// HighCyan returns the string representation of the given value in high-intensity cyan.
func HighCyan(v interface{}) string {
	return text.FgHiCyan.Sprint(v)
}

// HighGreen returns the string representation of the given value in high-intensity green.
func HighGreen(v interface{}) string {
	return text.FgHiGreen.Sprint(v)
}

// HighMagenta returns the string representation of the given value in high-intensity magenta.
func HighMagenta(v interface{}) string {
	return text.FgHiMagenta.Sprint(v)
}

// HighRed returns the string representation of the given value in high-intensity red.
func HighRed(v interface{}) string {
	return text.FgHiRed.Sprint(v)
}

// HighYellow returns the string representation of the given value in high-intensity yellow.
func HighYellow(v interface{}) string {
	return text.FgHiYellow.Sprint(v)
}

// HighWhite returns the string representation of the given value in high-intensity white.
func HighWhite(v interface{}) string {
	return text.FgHiWhite.Sprint(v)
}

// Blue returns the string representation of the given value in blue.
func Blue(v interface{}) string {
	return text.FgBlue.Sprint(v)
}

// Cyan returns the string representation of the given value in cyan.
func Cyan(v interface{}) string {
	return text.FgCyan.Sprint(v)
}

// Green returns the string representation of the given value in green.
func Green(v interface{}) string {
	return text.FgGreen.Sprint(v)
}

// Magenta returns the string representation of the given value in magenta.
func Magenta(v interface{}) string {
	return text.FgMagenta.Sprint(v)
}

// Red returns the string representation of the given value in red.
func Red(v interface{}) string {
	return text.FgRed.Sprint(v)
}

// Yellow returns the string representation of the given value in yellow.
func Yellow(v interface{}) string {
	return text.FgYellow.Sprint(v)
}

// White returns the string representation of the given value in white.
func White(v interface{}) string {
	return text.FgWhite.Sprint(v)
}
