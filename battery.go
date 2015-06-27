package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func getBatteryStatus(path string) (status string, err error) {
	present, err := ioutil.ReadFile(fmt.Sprintf("%s/present", path))
	if string(present) != "1" {
		status = "No battery"
	}

	energy_now, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_now", path))
	if err != nil {
		return
	}
	energy_full, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_full", path))
	if err != nil {
		return
	}
	charging, err := ioutil.ReadFile(fmt.Sprintf("%s/status", path))
	if err != nil {
		return
	}

	var ch string
	c := strings.TrimSpace(string(charging))
	switch c {
	case "Charging":
		ch = "+"
	case "Discharging":
		ch = "-"
	case "Full":
		ch = "="
	default: //Something went wrong determining if we are charging or discharging.
		ch = "?"
	}
	var enow, efull int
	fmt.Sscanf(string(energy_now), "%d", &enow)
	fmt.Sscanf(string(energy_full), "%d", &efull)

	//Format the status message
	status = fmt.Sprintf("Bat: %s %d%%", ch, enow*100/efull)
	return
}
