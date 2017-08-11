package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/tarm/goserial"
)

var conf config

const pulsesPerLiterForFt330 int = 2724

func main() {
	configFile := flag.String("config", "kaskpi.toml", "configuration file")
	flag.Parse()

	parseConfig(*configFile)

	// socat -d -d pty,raw,echo=0 stdio
	c := &serial.Config{Name: conf.Arduino.SerialPort, Baud: conf.Arduino.BaudRate}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	sendCommand(port, VersionCommand{})

	var pins []int
	for key, _ := range conf.Arduino.PinToTap {
		keyAsInt, _ := strconv.Atoi(key)
		pins = append(pins, keyAsInt)
	}
	sendCommand(port, Ft330InitializeCommand{PinCount: len(pins), Pins: pins, Delay: conf.Arduino.Delay, Threshold: conf.Arduino.Threshold})
	sendCommand(port, WiegandInitalizeCommand{Data0Pin: conf.Arduino.WiegandPins["data0"], Data1Pin: conf.Arduino.WiegandPins["data1"]})

	eventEmitter := getEmitter()
	waitForEvents(port, eventEmitter)
}
