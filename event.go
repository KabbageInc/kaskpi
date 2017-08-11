package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func waitForEvents(port io.ReadWriteCloser, eventEmitter io.Writer) {
	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		resp := scanner.Text()

		if len(resp) > 0 {
			dispatchEvent(eventEmitter, resp)
		}
	}
}

func dispatchEvent(eventEmitter io.Writer, event string) {
	parts := strings.Split(event, ":")

	if len(parts) < 2 {
		fmt.Println("invalid event: " + event)
		return
	}

	var eventName = parts[0]
	var timestampUint64, timestampError = strconv.ParseUint(parts[1], 10, 64)

	if timestampError != nil {
		fmt.Println("invalid event timestamp: " + event)
		return
	}
	var timestamp = uint32(timestampUint64)

	switch eventName {
	case "heartbeat":
		processHeartbeat(eventEmitter, timestamp)
	case "ft330_start":
		processFt330PourStart(eventEmitter, timestamp, parts[2:])
	case "ft330_end":
		processFt330PourEnd(eventEmitter, timestamp, parts[2:])
	case "wiegand_state":
		processWiegandState(eventEmitter, timestamp, parts[2:])
	case "wiegand_receive":
		processWiegandReceive(eventEmitter, timestamp, parts[2:])
	default:
		fmt.Println("unknown event: " + event)
	}
}

func processFt330PourStart(eventEmitter io.Writer, timestamp uint32, eventPayload []string) {
	fmt.Printf("FT-330 pour start: timestamp=%v, payload=%v", timestamp, strings.Join(eventPayload, ":"))
	fmt.Println()

	pinRaw := eventPayload[0]
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

	eventEmitter.Write(serializeMessage(msg))
}

func processFt330PourEnd(eventEmitter io.Writer, timestamp uint32, eventPayload []string) {
	fmt.Printf("FT-330 pour end: timestamp=%v, payload=%v", timestamp, strings.Join(eventPayload, ":"))
	fmt.Println()

	pinRaw := eventPayload[0]
	pin, err := strconv.Atoi(pinRaw)

	if err != nil {
		fmt.Println("invalid pin: " + pinRaw)
		fmt.Println(err.Error())
		return
	}

	pulsesRaw := eventPayload[1]
	pulses, err := strconv.Atoi(pulsesRaw)

	if err != nil {
		fmt.Println("invalid pulses: " + pinRaw)
		fmt.Println(err.Error())
		return
	}

	durationRaw := eventPayload[2]
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

	eventEmitter.Write(serializeMessage(msg))
}

func processWiegandState(eventEmitter io.Writer, timestamp uint32, eventPayload []string) {
	fmt.Printf("Wiegand state: timestamp=%v, payload=%v", timestamp, strings.Join(eventPayload, ":"))
	fmt.Println()

	connected, err := strconv.ParseBool(eventPayload[0])

	if err != nil {
		fmt.Println("invalid state: " + eventPayload[0])
		fmt.Println(err.Error())
		return
	}

	msg := WiegandStateMessage{
		Message:   Message{EventType: "PourStart", Timestamp: time.Now()},
		Connected: connected,
	}

	eventEmitter.Write(serializeMessage(msg))
}

func processWiegandReceive(eventEmitter io.Writer, timestamp uint32, eventPayload []string) {
	fmt.Printf("Wiegand receive: timestamp=%v, payload=%v", timestamp, strings.Join(eventPayload, ":"))
	fmt.Println()

	bitLengthRaw := eventPayload[0]
	bitLength, err := strconv.Atoi(bitLengthRaw)

	if err != nil {
		fmt.Println("invalid pin: " + bitLengthRaw)
		fmt.Println(err.Error())
		return
	}

	code := eventPayload[1]

	msg := WiegandReceiveMessage{
		Message:   Message{EventType: "PourStart", Timestamp: time.Now()},
		BitLength: bitLength,
		Code:      code,
	}

	eventEmitter.Write(serializeMessage(msg))
}

func processHeartbeat(eventEmitter io.Writer, timestamp uint32) {
	fmt.Printf("heartbeat: timestamp=%v", timestamp)
	fmt.Println()

	//msg := HeartbeatMessage{Message: Message{EventType: "Heartbeat", Timestamp: time.Now()}}
	//eventEmitter.Write(serializeMessage(msg))
}
