package pkg

import "bytes"

var magicBytes = map[string][]byte{
	"pdf":  {0x25, 0x50, 0x44, 0x46},
	"doc":  {0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1},
	"docx": {0x50, 0x4B, 0x03, 0x04},
}

func ValidateFileType(ext string, header []byte) bool {
	expected, ok := magicBytes[ext]
	if !ok {
		return false
	}
	return bytes.HasPrefix(header, expected)
}

func IsValidResumeType(fileType string) bool {
	validTypes := map[string]bool{"pdf": true, "doc": true, "docx": true}
	return validTypes[fileType]
}