package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Command interface {
	AsString() string
}

type VersionCommand struct {
	Command
}

func (self VersionCommand) AsString() string {
	return "version"
}

type Ft330InitializeCommand struct {
	Command
	PinCount  int
	Pins      []int
	Delay     int
	Threshold int
}

func (self Ft330InitializeCommand) AsString() string {
	return fmt.Sprintf("ft330_init:%v:%v:%v:%v", self.PinCount, SplitToString(self.Pins, ","), self.Delay, self.Threshold)
}

// this is dumb...
func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}

	return strings.Join(b, sep)
}

type WiegandInitalizeCommand struct {
	Command
	Data0Pin int
	Data1Pin int
}

func (self WiegandInitalizeCommand) AsString() string {
	return fmt.Sprintf("wiegand_init:%v:%v", self.Data0Pin, self.Data1Pin)
}

func sendCommand(port io.ReadWriteCloser, command Command) {
	sendCommandString(port, command.AsString())
}

func sendCommandString(port io.ReadWriteCloser, cmd string) {
	fmt.Println("sending command: " + cmd)

	port.Write([]byte(cmd))
	port.Write([]byte("\r\n"))

	scanner := bufio.NewScanner(port)

	for scanner.Scan() {
		resp := scanner.Text()

		if len(resp) > 0 {
			fmt.Println("response: " + resp)
			break
		}
	}
}
