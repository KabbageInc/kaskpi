package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Name    string
	Payload []string
}

func waitForEvents(port io.ReadWriteCloser) {
	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		resp := scanner.Text()

		if len(resp) > 0 {
			if event, err := parseEvent(resp); err == nil {
				event.dispatch()
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func parseEvent(eventString string) (Event, error) {
	parts := strings.Split(eventString, ":")
	event := Event{}
	if len(parts) < 2 {
		return event, errors.New("invalid event: " + eventString)
	}

	event.Name = parts[0]
	event.Payload = parts[2:]
	return event, nil
}

func (event Event) dispatch() {
	switch event.Name {
	case "heartbeat":
		event.processHeartbeat()
	case "ft330_start":
		event.processFt330PourStart()
	case "ft330_end":
		event.processFt330PourEnd()
	case "wiegand_state":
		event.processWiegandState()
	case "wiegand_receive":
		event.processWiegandReceive()
	default:
		fmt.Println("unknown event: " + event.Name)
	}
}

func (event Event) processFt330PourStart() {
	fmt.Printf("FT-330 pour start: timestamp=%v, payload=%v", time.Now().Unix(), strings.Join(event.Payload, ":"))
	fmt.Println()

	pinRaw := event.Payload[0]
	pin, err := strconv.Atoi(pinRaw)

	if err != nil {
		fmt.Println("invalid pin: " + pinRaw)
		fmt.Println(err.Error())
		return
	}

	msg := PourStartMessage{
		TapMessage: TapMessage{
			Message: Message{
				EventType: "PourStart",
				Timestamp: time.Now()},
			Tap: pinToTap(pin)}}

	getEmitter().Write(serializeMessage(msg))
}

func (event Event) processFt330PourEnd() {
	fmt.Printf("FT-330 pour end: timestamp=%v, payload=%v\n", time.Now().Unix(), strings.Join(event.Payload, ":"))

	pinRaw := event.Payload[0]
	pin, err := strconv.Atoi(pinRaw)

	if err != nil {
		fmt.Println("invalid pin: " + pinRaw)
		fmt.Println(err.Error())
		return
	}

	pulsesRaw := event.Payload[1]
	pulses, err := strconv.Atoi(pulsesRaw)

	if err != nil {
		fmt.Println("invalid pulses: " + pinRaw)
		fmt.Println(err.Error())
		return
	}

	durationRaw := event.Payload[2]
	duration, err := strconv.Atoi(durationRaw)

	if err != nil {
		fmt.Println("invalid duration: " + durationRaw)
		fmt.Println(err.Error())
		return
	}

	msg := PourEndMessage{
		TapMessage:         TapMessage{Message: Message{EventType: "PourEnd", Timestamp: time.Now()}, Tap: pinToTap(pin)},
		Milliliters:        float64(pulses) / float64(pulsesPerLiterForFt330) * 1000,
		Duration:           duration,
		RawFt330SensorData: RawFt330SensorData{Pulses: pulses}}

	getEmitter().Write(serializeMessage(msg))
}

func (event Event) processWiegandState() {
	fmt.Printf("Wiegand state: timestamp=%v, payload=%v\n", time.Now().Unix(), strings.Join(event.Payload, ":"))

	connected, err := strconv.ParseBool(event.Payload[0])

	if err != nil {
		fmt.Println("invalid state: " + event.Payload[0])
		fmt.Println(err.Error())
		return
	}

	msg := WiegandStateMessage{
		Message:   Message{EventType: "PourStart", Timestamp: time.Now()},
		Connected: connected,
	}

	getEmitter().Write(serializeMessage(msg))
}

func (event Event) processWiegandReceive() {
	fmt.Printf("Wiegand receive: timestamp=%v, payload=%v\n", time.Now().Unix(), strings.Join(event.Payload, ":"))

	bitLengthRaw := event.Payload[0]
	bitLength, err := strconv.Atoi(bitLengthRaw)

	if err != nil {
		fmt.Println("invalid pin: " + bitLengthRaw)
		fmt.Println(err.Error())
		return
	}

	code := event.Payload[1]

	msg := WiegandReceiveMessage{
		Message:   Message{EventType: "PourStart", Timestamp: time.Now()},
		BitLength: bitLength,
		Code:      code,
	}

	getEmitter().Write(serializeMessage(msg))
}

func (event Event) processHeartbeat() {
	fmt.Printf("heartbeat: timestamp=%v\n", time.Now().Unix())

	msg := HeartbeatMessage{Message: Message{EventType: "Heartbeat", Timestamp: time.Now()}}
	getEmitter().Write(serializeMessage(msg))
}
