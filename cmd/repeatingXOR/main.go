package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/racsoraul/cryptopals/set/one"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var filePath string
	var decryptMode bool
	flag.StringVar(&filePath, "f", "", "Path to file to encrypt/decrypt.")
	flag.BoolVar(&decryptMode, "d", false, "Indicates if it should encrypt or decrypt the file.")
	flag.Parse()
	if len(filePath) == 0 {
		log.Fatalln("Path to file -f is required.")
	}

	fmt.Print("Enter Key: ")
	scanner := bufio.NewScanner(os.Stdin)
	var key string
	for scanner.Scan() {
		key = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	inputFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()

	tmpFileName := fmt.Sprintf("%d-tmp-%s", time.Now().Unix(), filepath.Base(inputFile.Name()))
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer tmpFile.Close()

	inputScanner := bufio.NewScanner(inputFile)

	if decryptMode {
		for inputScanner.Scan() {
			text, err := one.DecodeHex(inputScanner.Text())
			if err != nil {
				log.Fatalln(err)
			}
			_, err = tmpFile.WriteString(fmt.Sprintf("%s\n", one.EncryptWithRepeatingXOR(string(text), key)))
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		for inputScanner.Scan() {
			text := inputScanner.Text()
			_, err := tmpFile.WriteString(fmt.Sprintf("%s\n", one.EncodeToHex(one.EncryptWithRepeatingXOR(text, key))))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	err = inputScanner.Err()
	if err != nil {
		log.Fatalln(err)
	}

	// Remove original file.
	err = os.Remove(inputFile.Name())
	if err != nil {
		log.Fatalln(err)
	}

	// Replace by encrypted/decrypted version.
	err = os.Rename(tmpFile.Name(), inputFile.Name())
	if err != nil {
		log.Fatalln(err)
	}
}
