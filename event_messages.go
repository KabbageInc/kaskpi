package main

import "time"

type Message struct {
	EventType string
	Timestamp time.Time
}

type TapMessage struct {
	Message
	Tap int
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

type WiegandStateMessage struct {
	Message
	Connected bool
}

type WiegandReceiveMessage struct {
	Message
	BitLength int
	Code      string
}
