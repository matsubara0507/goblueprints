package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	OtherWord  string   `yaml:"otherWord"`
	Transforms []string `yaml:"transforms"`
}

func main() {
	var filename = flag.String("config", "config.yaml", "変換方法")
	flag.Parse()

	buf, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Fatal("transforms file is not found.")
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatal("error Unmarshal: ", buf)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		t := config.Transforms[rand.Intn(len(config.Transforms))]
		fmt.Println(strings.Replace(t, config.OtherWord, s.Text(), -1))
	}

	err = s.Err()
	if err != nil {
		log.Fatal("error: ", err)
	}
}
