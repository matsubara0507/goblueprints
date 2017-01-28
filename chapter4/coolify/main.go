package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func duplicateVowel(word []byte, i int) []byte {
	return append(word[:i+1], word[i:]...)
}

func removeVowel(word []byte, i int) []byte {
	return append(word[:i], word[i+1:]...)
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := []byte(s.Text())
		if randBool() {
			var vowelIndexes []int
			for i, char := range word {
				switch char {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					vowelIndexes = append(vowelIndexes, i)
				}
			}
			if len(vowelIndexes) > 0 {
				vI := vowelIndexes[rand.Intn(len(vowelIndexes))]
				if randBool() {
					word = duplicateVowel(word, vI)
				} else {
					word = removeVowel(word, vI)
				}
			}
		}
		fmt.Println(string(word))
	}

	err := s.Err()
	if err != nil {
		log.Fatal("error: ", err)
	}
}
