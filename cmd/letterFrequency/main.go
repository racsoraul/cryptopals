package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type letterFrequency map[byte]int
type letterDistribution struct {
	letter byte
	value  float32
}

func main() {
	var (
		fileName        string
		letters         uint
		customLetterSet string
	)

	flag.StringVar(&fileName, "file", "", "Path to file to process.")
	flag.UintVar(&letters, "letters", 0, "Set of letters to find frequency. Values: 0: a-z, 1: A-Z, 2: a-zA-Z.")
	flag.StringVar(&customLetterSet, "custom", "", "Custom set of letters to find frequency.")

	flag.Parse()

	if len(fileName) == 0 {
		log.Println("'file' argument is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var letterSet []byte
	if len(customLetterSet) > 0 {
		letterSet = []byte(customLetterSet)
	} else {
		switch letters {
		case 0:
			letterSet = []byte("abcdefghijklmnopqrstuvwxyz")
		case 1:
			letterSet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		case 2:
			letterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		default:
			log.Fatalf("Value %d is invalid for argument 'letters'.\n", letters)
		}
	}

	fmt.Printf(
		"Running using values:\n -fileName: %s\n -letters: %d -> %s\n -custom: %s\n",
		fileName,
		letters, letterSet,
		customLetterSet,
	)

	textFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer textFile.Close()

	reader := bufio.NewReader(textFile)
	buffer := make([]byte, 64<<10)
	frequencyCounter := make(letterFrequency)
	for {
		n, errReader := reader.Read(buffer)
		if errReader == io.EOF {
			break
		}
		if errReader != nil {
			log.Println(errReader)
		}

		for _, letter := range letterSet {
			count := frequencyCounter[letter]
			frequencyCounter[letter] = count + bytes.Count(buffer[:n], []byte{letter})
		}
	}
	distribution := make([]letterDistribution, 0)
	var totalCount int
	for _, count := range frequencyCounter {
		totalCount += count
	}
	for key, count := range frequencyCounter {
		distribution = append(distribution, letterDistribution{
			letter: key,
			value:  float32(count) / float32(totalCount),
		})
	}
	sort.SliceStable(distribution, func(i, j int) bool {
		return distribution[i].letter < distribution[j].letter
	})
	for _, frequency := range distribution {
		fmt.Printf("%c: %.8f\n", frequency.letter, frequency.value)
	}
}
