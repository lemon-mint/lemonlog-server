package main

import (
	"time"

	"github.com/lemon-mint/lemonlog-server/backend/memory"
)

var db = memory.New()

func pusher() {
	for {
		logs := getQ(128)
		if len(logs) > 0 {
			db.Put(logs)
		} else {
			time.Sleep(time.Millisecond * 1)
		}
	}
}
