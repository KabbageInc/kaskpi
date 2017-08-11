package main

import (
	"encoding/json"
	"io"
	"time"
)

type Message struct {
	EventType string
	Timestamp time.Time
}

func (m Message) Write(w io.Writer) {
	jsonMessage, _ := json.Marshal(m)
	w.Write(jsonMessage)
}

type TapMessage struct {
	Message
	Tap int
}

func (m TapMessage) Write(w io.Writer) {
	jsonMessage, _ := json.Marshal(m)
	w.Write(jsonMessage)
}

type HeartbeatMessage struct {
	Message
}

type PourStartMessage struct {
	TapMessage
}

type RawFt330SensorData struct {
	Pulses int
}

type PourEndMessage struct {
	TapMessage
	Milliliters float64
	Duration    int
	RawFt330SensorData
}

func (m PourEndMessage) Write(w io.Writer) {
	jsonMessage, _ := json.Marshal(m)
	w.Write(jsonMessage)
}

type WiegandStateMessage struct {
	Message
	Connected bool
}

func (m WiegandStateMessage) Write(w io.Writer) {
	jsonMessage, _ := json.Marshal(m)
	w.Write(jsonMessage)
}

type WiegandReceiveMessage struct {
	Message
	BitLength int
	Code      string
}

func (m WiegandReceiveMessage) Write(w io.Writer) {
	jsonMessage, _ := json.Marshal(m)
	w.Write(jsonMessage)
}
