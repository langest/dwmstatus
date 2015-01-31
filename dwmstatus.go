package main

/*
#cgo LDFLAGS: -lX11 -lasound
#include <X11/Xlib.h>
#include "getvol.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

const (
	mpdAddr = "localhost:6600"
)

var (
	dpy = C.XOpenDisplay(nil)
)

func getVolumePerc() int {
	return int(C.get_volume_perc())
}

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

func getLoadAverage(file string) (lavg string, err error) {
	loadavg, err := ioutil.ReadFile(file)
	if err != nil {
		return "Couldn't read loadavg", err
	}
	lavg = strings.Join(strings.Fields(string(loadavg))[:3], " ")
	return
}

func setStatus(s *C.char) {
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), s)
	C.XSync(dpy, 1)
}

func nowPlaying(addr string) (np string, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		np = "Couldn't connect to mpd."
		return
	}
	defer conn.Close()
	reply := make([]byte, 512)
	conn.Read(reply) // The mpd OK has to be read before we can actually do things.

	message := "status\n"
	conn.Write([]byte(message))
	conn.Read(reply)
	r := string(reply)
	arr := strings.Split(string(r), "\n")
	if arr[8] != "state: play" { //arr[8] is the state according to the mpd documentation
		status := strings.SplitN(arr[8], ": ", 2)
		np = fmt.Sprintf("mpd - [%s]", status[1]) //status[1] should now be stopped or paused
		return
	}

	message = "currentsong\n"
	conn.Write([]byte(message))
	conn.Read(reply)
	r = string(reply)
	arr = strings.Split(string(r), "\n")
	if len(arr) > 5 {
		artist := ""
		title := ""
		for _, info := range arr {
			field := strings.SplitN(info, ":", 2)
			switch field[0] {
			case "Artist":
				artist = field[1]
			case "Title":
				title = field[1]
			default:
				//do nothing with the field
			}
		}
		np = artist + " - " + title
		return
	} else { //This is a nonfatal error.
		np = "Playlist is empty."
		return
	}
}

func formatStatus(format string, args ...interface{}) *C.char {
	status := fmt.Sprintf(format, args...)
	return C.CString(status)
}

func main() {
	if dpy == nil {
		log.Fatal("Can't open display")
	}
	for {
		tim := time.Now().Format("Mon 02 Jan 15:04")
		bat, _ := getBatteryStatus("/sys/class/power_supply/BAT0")
		//if err != nil {
		//	log.Println(err)
		//}
		//lavg, _ := getLoadAverage("/proc/loadavg")
		//		if err != nil {
		//			log.Println(err)
		//		}
		mpd, _ := nowPlaying(mpdAddr)
		//		if err != nil {
		//			log.Println(err)
		//		}
		vol := getVolumePerc()
		s := formatStatus("%s || Vol: %d%% || %s || %s", mpd, vol, bat, tim)
		setStatus(s)
		time.Sleep(time.Second * 1)
	}
}
