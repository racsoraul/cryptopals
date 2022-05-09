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
	flag.StringVar(&filePath, "f", "", "Path to file to encrypt.")
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

	for inputScanner.Scan() {
		text := inputScanner.Text()
		_, err := tmpFile.WriteString(one.EncryptWithRepeatingXOR(text, key) + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}

	err = inputScanner.Err()
	if err != nil {
		log.Fatalln(err)
	}
}
