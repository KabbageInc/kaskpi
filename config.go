package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strconv"
	"time"
)

type config struct {
	Aws     awsConfig
	Arduino arduinoConfig
}

type awsConfig struct {
	Region            string
	SnsArn            string   `toml:"sns_arn"`
	HeartbeatInterval duration `toml:"heartbeat_interval"`
	AccessKeyId       string   `toml:"access_key_id"`
	SecretAccessKey   string   `toml:"secret_access_key"`
}

type arduinoConfig struct {
	SerialPort  string         `toml:"serial_port"`
	BaudRate    int            `toml:"baud_rate"`
	PinToTap    map[string]int `toml:"pin_to_tap"`
	Delay       int
	Threshold   int
	WiegandPins map[string]int `toml:"wiegand_pins"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func parseConfig(configFile string) {
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatal("config file is missing: ", configFile)
	}

	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		log.Fatal("failed to parse config: " + err.Error())
	}

	fmt.Printf("config: %+v\n", conf)
	fmt.Println()
}

func pinToTap(pin int) int {
	return conf.Arduino.PinToTap[strconv.Itoa(pin)]
}
