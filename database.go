package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

var eventDb *bolt.DB
var bucket = []byte("eventbucket")

func initDb(path string) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		//todo set no db mode?
		log.Fatal(err.Error())
	}
	eventDb = db
	err = eventDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("error create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

func closeDb() {
	eventDb.Close()
}

func writeToDb(p []byte) {
	err := eventDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		id, _ := b.NextSequence()
		return b.Put(itob(id), p)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func consumeNext() (bool, []byte) {
	var k, v []byte
	err := eventDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		k, v = b.Cursor().First()
		if k != nil {
			return b.Delete(k)
		}
		return nil
	})
	return err == nil && k != nil, v
}

func getStoredEventCount() (int, error) {
	var count int
	err := eventDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}
