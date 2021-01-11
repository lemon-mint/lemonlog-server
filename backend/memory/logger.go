package memory

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/allegro/bigcache"
	"github.com/lemon-mint/lemonlog-server/backend"
)

//MemLogger :
type MemLogger struct {
	db *bigcache.BigCache
}

//New database
func New() backend.LogStore {
	var database, _ = bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 720))
	var logger backend.LogStore = MemLogger{database}
	return logger
}

//Put to in memory database
func (db MemLogger) Put(data []*backend.Log) error {
	var buf bytes.Buffer
	for i := range data {
		buf.Reset()
		encoder := gob.NewEncoder(&buf)
		encoder.Encode(*data[i])
		db.db.Set(data[i].UUID, buf.Bytes())
	}
	return nil
}

//Del : Delete log
func (db MemLogger) Del(uuid string) error {
	return db.db.Delete(uuid)
}

//Get : Get log
func (db MemLogger) Get(uuid string) (*backend.Log, error) {
	data, err := db.db.Get(uuid)
	if err != nil {
		return &backend.Log{}, err
	}
	var buf bytes.Buffer
	var log backend.Log
	buf.Write(data)
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&log)
	if err != nil {
		return &backend.Log{}, err
	}
	return &log, nil
}
