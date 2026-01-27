package crypto

import (
    "encoding/base64"
)

// EncodeBase64 encodes bytes to base64 string
func EncodeBase64(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes base64 string to bytes
func DecodeBase64(encoded string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(encoded)
}