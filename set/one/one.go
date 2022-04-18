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

var EnglishLettersDistribution = map[byte]float64{
	'a': 0.07743208627550165,
	'b': 0.01402241586697527,
	'c': 0.02665670667329359,
	'd': 0.04920785702311875,
	'e': 0.13464518994079883,
	'f': 0.025036247121552113,
	'g': 0.017007472935972733,
	'h': 0.05719839895067157,
	'i': 0.06294794236928244,
	'j': 0.001267546400727001,
	'k': 0.005084890317533608,
	'l': 0.03706176274237046,
	'm': 0.030277007414117114,
	'n': 0.07125316518982316,
	'o': 0.07380002176297765,
	'p': 0.017513315119093483,
	'q': 0.0009499245648139707,
	'r': 0.06107162078305546,
	's': 0.061262782073188304,
	't': 0.08760480785349399,
	'u': 0.030426995503298266,
	'v': 0.01113735085743191,
	'w': 0.02168063124398945,
	'x': 0.0019880774173815607,
	'y': 0.022836421813561863,
	'z': 0.0006293617859758195,
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
	cipherText, err := decodeHex(src)
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
