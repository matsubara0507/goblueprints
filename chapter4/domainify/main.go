package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

var tlds = []string{"com", "net"}

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_-"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var tdl = flag.String("tdl", tlds[rand.Intn(len(tlds))], "トップレベルドメイン")
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(s.Text())
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if strings.ContainsRune(allowedChars, r) {
				newText = append(newText, r)
			}
		}
		fmt.Println(string(newText) + "." + *tdl)
	}

	err := s.Err()
	if err != nil {
		log.Fatal("error: ", err)
	}
}
