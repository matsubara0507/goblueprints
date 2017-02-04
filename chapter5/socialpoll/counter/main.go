package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const updateDuration = 1 * time.Second

func main() {
	err := counterMain()
	if err != nil {
		log.Fatal(err)
	}
}

func counterMain() error {
	log.Println("データベースに接続しています...")
	db, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}
	defer func() {
		log.Println("データベース接続を閉じます...")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	var countsLock sync.Mutex
	var counts map[string]int

	log.Println("NSQ に接続します...")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		return err
	}

	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))

	err = q.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		return err
	}

	log.Println("NSQ 上での投票を待機します...")
	ticker := time.NewTicker(updateDuration)
	defer ticker.Stop()

	update := func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("新しい投票はありません。データベースの更新をスキップします")
			return
		}
		log.Println("データベースを更新します...")
		log.Println(counts)
		ok := true
		for option, count := range counts {
			sel := bson.M{"options": bson.M{"$in": []string{option}}}
			up := bson.M{"$inc": bson.M{"result" + option: count}}
			_, err := pollData.UpdateAll(sel, up)
			if err != nil {
				log.Println("更新に失敗しました: ", err)
				ok = false
			} else {
				counts[option] = 0
			}
		}
		if ok {
			log.Println("データベースの更新が完了しました")
			counts = nil
		}
	}

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		select {
		case <-ticker.C:
			update()
		case <-termChan:
			q.Stop()
		case <-q.StopChan:
			return nil
		}
	}
}
