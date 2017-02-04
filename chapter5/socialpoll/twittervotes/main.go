package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	nsq "github.com/nsqio/go-nsq"

	mgo "gopkg.in/mgo.v2"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("MongoDB にダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("データベース接続が閉じられました")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) {
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	for vote := range votes {
		pub.Publish("votes", []byte(vote))
	}
	log.Println("Publisher: 停止中です")
	pub.Stop()
	log.Println("Publisher: 停止しました")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		cancel()
		log.Println("停止します...")
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	err := dialdb()
	if err != nil {
		log.Fatalln("MongoDB へのダイヤルに失敗しました: ", err)
	}
	defer closedb()

	// Start Processes
	votes := make(chan string)
	go twitterStream(ctx, votes)
	publishVotes(votes)
}
