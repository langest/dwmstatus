package main

/*
#cgo LDFLAGS: -lX11 -lasound
#include <X11/Xlib.h>
#include "getvol.h"
*/
import "C"

import (
	"fmt"
	"log"
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

func setStatus(s *C.char) {
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), s)
	C.XSync(dpy, 1)
}

func formatStatus(format string, args ...interface{}) *C.char {
	status := fmt.Sprintf(format, args...)
	return C.CString(status)
}

func main() {
	defer C.XCloseDisplay(dpy)
	if dpy == nil {
		log.Fatal("Can't open display")
	}

	for {
		tim := time.Now().Format("Mon 02 Jan 15:04")
		bat, err := getBatteryStatus("/sys/class/power_supply/BAT0")
		//if err != nil {
		//log.Println(err)
		//}
		//lavg, _ := getLoadAverage("/proc/loadavg")
		//		if err != nil {
		//			log.Println(err)
		//		}
		mpd, err := nowPlaying(mpdAddr)
		if err != nil {
			log.Println(err)
		}
		vol := getVolumePerc()
		s := formatStatus("%s || Vol: %d%% || %s || %s", mpd, vol, bat, tim)
		setStatus(s)
		time.Sleep(time.Second * 1)
	}
}
