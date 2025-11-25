package base

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// IsTextFile checks if a file is text, supporting UTF-8, UTF-16 (with BOM),
// and other standard text formats.
func IsTextFile(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		slog.Error("error opening file", "path", filename, "error", err)
		return false, err
	}
	defer f.Close()
	slog.Debug("file opened", "path", filename)

	// Read up to 512 bytes
	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		slog.Error("error reading up to 512 bytes", "error", err)
		return false, err
	}

	// empty files are treated as text
	if n == 0 {
		slog.Warn("empty file", "path", filename)
		return true, nil
	}

	slog.Debug("successfully read bytaes from file", "path", filename, "read", n)
	return IsText(buffer[:n])
}

func IsText(buffer []byte) (bool, error) {

	n := len(buffer)

	// 2. CHECK FOR BOM (Byte Order Mark)
	// This must happen BEFORE the NUL check, because UTF-16 contains NUL bytes.
	// If we find a BOM, we are 100% confident it is text.

	// Check for UTF-16 LE (Little Endian) - Hex: FF FE
	if n >= 2 && buffer[0] == 0xFF && buffer[1] == 0xFE {
		return true, nil
	}

	// Check for UTF-16 BE (Big Endian) - Hex: FE FF
	if n >= 2 && buffer[0] == 0xFE && buffer[1] == 0xFF {
		return true, nil
	}

	// Check for UTF-8 BOM - Hex: EF BB BF
	// (Note: The NUL check below handles UTF-8 fine, but this is a quick optimization)
	if n >= 3 && buffer[0] == 0xEF && buffer[1] == 0xBB && buffer[2] == 0xBF {
		return true, nil
	}

	// 3. CHECK FOR NUL BYTES
	// If it wasn't a BOM-marked UTF-16 file, and it contains a NUL byte,
	// it is almost certainly binary (like an image or compiled app).
	if bytes.IndexByte(buffer, 0) != -1 {
		return false, nil
	}

	// 4. MIME TYPE SNIFFING
	// Use net/http to guess the type based on the content signature.
	contentType := http.DetectContentType(buffer)

	// If it explicitly says "text/", it's text.
	if strings.HasPrefix(contentType, "text/") {
		return true, nil
	}

	// 5. WHITELIST "APPLICATION" TEXT TYPES
	whitelistedTypes := []string{
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-javascript",
		"image/svg+xml",
		"application/x-yaml",
	}

	for _, valid := range whitelistedTypes {
		if strings.HasPrefix(contentType, valid) {
			return true, nil
		}
	}

	// 6. FALLBACK
	// If the standard library sees "application/octet-stream" (unknown),
	// BUT we have passed the NUL check (Step 3), we assume it's text.
	// This catches README files, source code, logs, etc.
	if contentType == "application/octet-stream" {
		return true, nil
	}

	// Default to binary if it identified as something else (e.g., image/png)
	return false, nil
}

/*
func main() {
	// Example usage
	files := []string{"main.go", "example.png", "utf16-file.txt"}

	for _, file := range files {
		// Create dummy files for demonstration if they don't exist
		// (In a real app, you would just check existing files)
		if file == "utf16-file.txt" {
			createDummyUTF16(file)
		}

		isText, _ := IsTextFile(file)
		if isText {
			fmt.Printf("[%s] is TEXT\n", file)
		} else {
			fmt.Printf("[%s] is BINARY\n", file)
		}
	}
}

// Helper to create a fake UTF-16 file for testing
func createDummyUTF16(name string) {
	f, _ := os.Create(name)
	defer f.Close()
	// Write BOM (FF FE) + 'A' (41 00)
	f.Write([]byte{0xFF, 0xFE, 0x41, 0x00})
}
*/
