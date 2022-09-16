package one

import (
	"os"
	"testing"
)

// Challenge 1
func TestHexToBase64(t *testing.T) {
	b64, err := hexToBase64("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
	if err != nil {
		t.Fatal(err)
	}
	wants := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	if b64 != wants {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, b64)
	}
}

// Challenge 2
func TestFixedXOR(t *testing.T) {
	actual, err := fixedXOR("1c0111001f010100061a024b53535009181c", "686974207468652062756c6c277320657965")
	if err != nil {
		t.Fatal(err)
	}
	wants := "746865206b696420646f6e277420706c6179"
	if wants != string(EncodeToHex(actual)) {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, actual)
	}
}

// Challenge 3
func TestDecipherSingleByteXOR(t *testing.T) {
	msg, err := decipherSingleByteXOR("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	if err != nil {
		t.Fatal(err)
	}
	wants := "Cooking MC's like a pound of bacon"
	if msg.plainText != wants {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, msg.plainText)
	}
}

// Challenge 4
func TestDecipherSingleByteXORFromFile(t *testing.T) {
	file, err := os.Open("challenge4_data.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	msg, err := decipherSingleByteXORFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	wants := "Now that the party is jumping"
	if msg.plainText != wants {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, msg.plainText)
	}
}

// Challenge 5
func TestEncryptWithRepeatingXOR(t *testing.T) {
	msg := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"
	actual := EncodeToHex(EncryptWithRepeatingXOR(msg, key))
	wants := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	if string(actual) != wants {
		t.Fatalf("\nWants:\n%s\nActual:\n%s\n", wants, actual)
	}
}

// Challenge 6
func TestHammingDistance(t *testing.T) {
	distance, err := HammingDistance([]byte("this is a test"), []byte("wokka wokka!!!"))
	if err != nil {
		t.Fatal(err)
	}
	wants := 37
	if distance != wants {
		t.Fatalf("wants: %d; got: %d", wants, distance)
	}
}
