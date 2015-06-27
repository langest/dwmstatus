package main

import (
	"fmt"
	mpd "github.com/fhs/gompd/mpd"
)

var (
	addr string
	cli  *mpd.Client
)

func dialMPD(address string) (err error) {
	addr = address
	cli, err = mpd.Dial("tcp", addr)
	return
}

func closeMPD() {
	if nil != cli.Ping() {
		cli.Close()
	}
}

func nowPlaying() (np string, err error) {
	if cli == nil || nil != cli.Ping() {
		cli, err = mpd.Dial("tcp", addr)
		if err != nil {
			return
		}
	}
	reply, err := cli.Status()
	if err != nil {
		return
	}
	if reply["state"] != "play" {
		np = fmt.Sprintf("mpd - [%s]", reply["state"]) //state should now be stopped or paused
		return
	}

	// If playing we want to print some song info
	reply, err = cli.CurrentSong()
	if err != nil {
		return
	}

	np = reply["Artist"] + " - " + reply["Title"]
	return
}
