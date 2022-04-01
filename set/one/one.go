package one

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// decodeHex Decodes Hex string.
func decodeHex(src string) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, []byte(src))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// encodeToHex Encodes string to Hex.
func encodeToHex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// hexToBase64 Convert hex to base64
func hexToBase64(src string) (string, error) {
	dst, err := decodeHex(src)
	if err != nil {
		return "", err
	}
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(dst)))
	base64.StdEncoding.Encode(b64, dst)
	return string(b64), nil
}

// fixedXOR takes two equal-length hex values and produces their XOR combination.
func fixedXOR(a, b string) (string, error) {
	if len(a) != len(b) {
		return "", fmt.Errorf("arguments must be of equal length")
	}
	decodedA, err := decodeHex(a)
	if err != nil {
		return "", err
	}
	decodedB, err := decodeHex(b)
	if err != nil {
		return "", err
	}
	result := make([]byte, len(decodedA))
	for i := 0; i < len(result); i++ {
		result[i] = decodedA[i] ^ decodedB[i]
	}
	return string(encodeToHex(result)), nil
}
