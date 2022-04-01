package setOne

import "testing"

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
	if wants != actual {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, actual)
	}
}
