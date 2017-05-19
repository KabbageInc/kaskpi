# Kabbage Kaskpi Overview

Raspberry Pi-based component of [Kabbage Kask](https://kask.kabbage.com) which performs the following functions:

  * Interface to Kaskduino running on an Arduino over serial
    * Tested on the Alamode shield
    * Supports GEMS Ft-330 sensors and conversions to metric
    * Supports Wiegand RFID data
  * Interface to a Dallas DS18B20 1-Wire temperature sensors
  * Send data to Amazon SNS for processing by the server-side component of Kask

# Hardware

  * Raspberry Pi 3
  * BCM2837 GPIO 18 wired to reboot Alamode
  * 3.3V vs 5V

## Raspberry Pi Pinout

| Pin | BCM   | Name   | Notes                                                          | 
| --- | ----- |------- | -------------------------------------------------------------- |
| 8   | BCM14 | TXD    | Serial TX                                                      |
| 9   | BCM15 | RXD    | Serial RX                                                      |
| 18  | BCM24 | n/a    | Alamode Arduino reset line, facilitates programming via serial |

# Configuration

Config files are in [TOML](https://github.com/toml-lang/toml).

```toml
[aws]
region = "us-east-1"
access_key_id = "xxx"
secret_access_key = "xxx"
sns_arn = "arn:aws:sns:us-east-1:xxx:kask-keg-event"
heartbeat_interval = "5s"

[arduino]
serial_port = "/dev/ttyAMA0"
baud_rate = 115200
pin_to_tap = { 4 = 1, 5 = 2, 6 = 3, 7 = 4 }
delay = 300
threshold = 100
rfid_pins = { data0 = 2, data1 = 3 }
```

# Toolchain

Trying to follow idiomatic Golang. Used JetBrains Gogland, add more here.

# SNS Events Message Schemas

General message schema has at least the following elements:

```json
{
	"EventType": "PourStart",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi"
}
```

## Startup

```json
```

## Heartbeat

Used to detect the keg going offline, sent every 5 seconds by default.

## Pour Start

```json
{
	"EventType": "PourStart",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"Tap": 1,
	"Sensor": {
		"Model": "GEMS FT-330"
	}
}
```

## Pour End

```json
{
	"EventType": "PourEnd",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"Tap": 1,
	"Pour": {
		"Milliliters": 837.5,
		"Timespan": 4.12
	},
	"RawData": {
		"Pulses": 3382
	},
	"Sensor": {
		"Model": "GEMS FT-330"
	}
}
```

## Kick

```json
{
	"EventType": "Kick",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"Tap": 1,
	"Sensor": {
		"Model": "GEMS FT-330"
	
}
```
## Wiegand Connection

```json
{
	"EventType": "WiegandConnection",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"Status": "Connected"
}
```

## Wiegand Scan

```json
{
	"EventType": "WiegandScan",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"Badge": {
		"Type": "HID",
		"BitLength": 26,
		"Code": 30675981993,
	},
}
```

## Temperature
```json
{
	"EventType": "Temperature",
	"Timestamp": "2017-04-28T18:25:43.511Z",
	"Source": "beerpi",
	"TemperatureInCelsius": 39.1,
	"Sensor": {
		"Model": "DS18B20",
		"Serial": "28-00000582cf56",
	}
}

