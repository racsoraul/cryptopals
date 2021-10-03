package setOne

import (
	"encoding/base64"
	"encoding/hex"
)

// challengeOne Convert hex to base64
func challengeOne(src string) (string, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, []byte(src))
	if err != nil {
		return "", err
	}
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(dst)))
	base64.StdEncoding.Encode(b64, dst)
	return string(b64), nil
}
