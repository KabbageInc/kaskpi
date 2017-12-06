package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestOpen(t *testing.T) {
	db, err := bolt.Open(tempfile(), 0600, nil)
	if err != nil {
		t.Fatal(err)
	} else if db == nil {
		t.Fatal("can't open db")
	}
}

func TestAdd(t *testing.T) {
	initDb(tempfile())
	defer closeDb()
	data := []byte("bananas")

	writeToDb(data)
	assertMessageCount(t, 1)

	writeToDb(data)
	assertMessageCount(t, 2)
}

func TestAddThenConsume(t *testing.T) {
	initDb(tempfile())
	defer closeDb()

	data := []byte("bananas")
	writeToDb(data)
	assertMessageCount(t, 1)

	if hasMessage, message := consumeNext(); hasMessage {
		if string(message) != string(data) {
			t.Fatalf("expected %s, got %s", string(data), string(message))
		}
	} else {
		t.Fatal("consume returned no event")
	}

	assertMessageCount(t, 0)
}

func TestFifo(t *testing.T) {
	initDb(tempfile())
	defer closeDb()

	first := []byte("should be first")
	second := []byte("should be second")

	writeToDb(first)
	writeToDb(second)

	_, v := consumeNext()
	if string(first) != string(v) {
		t.Logf("expected %s, got %s", string(first), string(v))
	}

	_, v = consumeNext()
	if string(second) != string(v) {
		t.Fatalf("expected %s, got %s", string(second), string(v))
	}
}

func assertMessageCount(t *testing.T, expected int) {
	if count, err := getStoredEventCount(); err != nil {
		t.Fatalf("Error getting stored count: %s", err.Error())
	} else {
		if count != expected {
			t.Fatalf("expected %v events, got %v", expected, count)
		}
	}
}

func tempfile() string {
	f, err := ioutil.TempFile("", "kaskpi-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
