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

type lettersFrequency map[byte]int
type lettersDistribution struct {
	letter byte
	value  float64
}

func main() {
	var (
		fileName        string
		set             uint
		customLetterSet string
	)

	flag.StringVar(&fileName, "file", "", "Path to file to process.")
	flag.UintVar(&set, "set", 0, "Set of letters to find frequency. Values: 0: a-z, 1: A-Z, 2: a-zA-Z.")
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
		switch set {
		case 0:
			letterSet = []byte("abcdefghijklmnopqrstuvwxyz")
		case 1:
			letterSet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		case 2:
			letterSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		default:
			log.Fatalf("Value %d is invalid for argument 'set'.\n", set)
		}
	}

	fmt.Printf(
		"Running using values:\n -fileName: %s\n -set: %d -> %s\n -custom: %s\n",
		fileName, set, letterSet, customLetterSet,
	)

	textFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer textFile.Close()

	reader := bufio.NewReader(textFile)
	buffer := make([]byte, 64<<10)
	frequencyCounter := make(lettersFrequency)
	var totalCount int
	for {
		n, errReader := reader.Read(buffer)
		if errReader == io.EOF {
			break
		}
		if errReader != nil {
			log.Println(errReader)
		}

		for _, letter := range letterSet {
			prevCount := frequencyCounter[letter]
			frequencyCount := bytes.Count(buffer[:n], []byte{letter})
			totalCount += frequencyCount
			frequencyCounter[letter] = prevCount + frequencyCount
		}
	}

	distribution := make([]lettersDistribution, 0)
	for key, count := range frequencyCounter {
		distribution = append(distribution, lettersDistribution{
			letter: key,
			value:  float64(count) / float64(totalCount),
		})
	}

	sort.SliceStable(distribution, func(i, j int) bool {
		return distribution[i].letter < distribution[j].letter
	})
	for _, frequency := range distribution {
		fmt.Printf("'%c': %g,\n", frequency.letter, frequency.value)
	}
}
