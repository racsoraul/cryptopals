package setOne

import "testing"

func TestChallengeOne(t *testing.T) {
	b64, err := challengeOne("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
	if err != nil {
		t.Fatal(err)
	}
	wants := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	if b64 != wants {
		t.Fatalf("\nWants: %s\nActual: %s\n", wants, b64)
	}
}
