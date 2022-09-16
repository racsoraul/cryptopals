package one

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strings"
)

//go:generate go run gen.go -file=frankenstein.txt

// DecodeHex Decodes Hex string.
func DecodeHex(src string) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, []byte(src))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// EncodeToHex Encodes string to Hex.
func EncodeToHex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// hexToBase64 Convert hex to base64
func hexToBase64(src string) (string, error) {
	dst, err := DecodeHex(src)
	if err != nil {
		return "", err
	}
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(dst)))
	base64.StdEncoding.Encode(b64, dst)
	return string(b64), nil
}

// fixedXOR takes two equal-length hex values and produces their XOR combination.
func fixedXOR(a, b string) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("arguments must be of equal length")
	}
	decodedA, err := DecodeHex(a)
	if err != nil {
		return nil, err
	}
	decodedB, err := DecodeHex(b)
	if err != nil {
		return nil, err
	}

	result := make([]byte, len(decodedA))
	for i := 0; i < len(result); i++ {
		result[i] = decodedA[i] ^ decodedB[i]
	}
	return result, nil
}

// scoreText Gives a score to text to help determine if it's English text.
// The closest to zero the better.
func scoreText(text []byte) float64 {
	var score float64
	lettersFrequency := make(map[byte]int)
	var totalCount int
	for letter := range EnglishLettersDistribution {
		count := bytes.Count(text, []byte{letter})
		totalCount += count
		lettersFrequency[letter] = count
	}

	lettersDistribution := make(map[byte]float64)
	for letter, frequency := range lettersFrequency {
		if totalCount > 0 {
			lettersDistribution[letter] = float64(frequency) / float64(totalCount)
		}
	}

	if len(lettersDistribution) == 0 {
		return math.Inf(1)
	}

	for letter, refDistribution := range EnglishLettersDistribution {
		distribution := lettersDistribution[letter]
		score += math.Abs(refDistribution - distribution)
	}

	return score
}

type guess struct {
	plainText string
	score     float64
	key       string
}

// decipherSingleByteXOR Deciphers the message.
func decipherSingleByteXOR(src string) (guess, error) {
	cipherText, err := DecodeHex(src)
	if err != nil {
		return guess{}, err
	}

	var candidateKeys [256]byte
	for i := range candidateKeys {
		candidateKeys[i] = byte(i)
	}

	candidatePlainText := make([]byte, len(cipherText))
	msgGuess := guess{score: math.Inf(1)}
	for _, key := range candidateKeys {
		for i := 0; i < len(candidatePlainText); i++ {
			candidatePlainText[i] = cipherText[i] ^ key
		}
		score := scoreText(candidatePlainText)
		if msgGuess.score > score {
			msgGuess.score = score
			msgGuess.plainText = string(candidatePlainText)
			msgGuess.key = string(key)
		}
	}
	return msgGuess, nil
}

// decipherSingleByteXORFromFile Decipher line containing the message from reader.
func decipherSingleByteXORFromFile(file *os.File) (guess, error) {
	scanner := bufio.NewScanner(file)

	msgGuess := guess{score: math.Inf(1)}
	for scanner.Scan() {
		bestGuess, err := decipherSingleByteXOR(scanner.Text())
		if err != nil {
			return guess{}, err
		}

		if bestGuess.score < msgGuess.score {
			msgGuess.score = bestGuess.score
			msgGuess.plainText = strings.TrimSpace(bestGuess.plainText)
		}
	}

	if err := scanner.Err(); err != nil {
		return guess{}, err
	}

	return msgGuess, nil
}

// EncryptWithRepeatingXOR Encrypts msg under the given key and returns the
// hex representation.
func EncryptWithRepeatingXOR(msg, key string) []byte {
	msgBytes := []byte(msg)
	keyBytes := []byte(key)
	encryptedMsg := make([]byte, len(msgBytes))

	keyLength := len(keyBytes)
	for i, msgByte := range msgBytes {
		encryptedMsg[i] = keyBytes[i%keyLength] ^ msgByte
	}

	return encryptedMsg
}

// HammingDistance Computes the edit distance between two strings.
func HammingDistance(a, b []byte) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("arguments must be of equal length")
	}

	var editDistance int
	for i := 0; i < len(a); i++ {
		distance := a[i] ^ b[i]
		for j := 0; j < 8; j++ {
			editDistance += int(distance & 1)
			distance = distance >> 1
		}
	}
	return editDistance, nil
}
