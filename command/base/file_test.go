package base

import (
	"bytes"
	"testing"
)

func TestIsTextFile(t *testing.T) {

	for filename, expected := range map[string]bool{
		"base.go":      true,
		"file.go":      true,
		"file_test.go": true,
		//"../../dist/archetype_linux_amd64_v1/archetype": false,
	} {
		result, err := IsTextFile(filename)
		if err != nil {
			t.Fatalf("cannot open file %q: %v", filename, err)
		}
		if result != expected {
			t.Fatalf("invalid detection for file %q: expected %t got %t", filename, expected, result)
		}
	}
}

func TestIsTextUTF16(t *testing.T) {
	var buffer bytes.Buffer
	buffer.Write([]byte{0xFF, 0xFE, 0x41, 0x00})
	if isText, _ := IsText(buffer.Bytes()); !isText {
		t.Fatalf("invalid detection for UTF-16 byte array: expected %t got %t", true, isText)
	}

}
