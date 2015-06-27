package main

import (
	"fmt"
	"net"
	"strings"
)

func nowPlaying(addr string) (np string, err error) {
	mpdConn, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}
	reply := make([]byte, 512)
	mpdConn.Read(reply) // The mpd OK has to be read before we can actually do things.

	defer mpdConn.Close()

	message := "status\n"
	mpdConn.Write([]byte(message))
	mpdConn.Read(reply)
	r := string(reply)
	arr := strings.Split(string(r), "\n")
	if arr[8] != "state: play" { //arr[8] is the state according to the mpd documentation
		status := strings.SplitN(arr[8], ": ", 2)
		np = fmt.Sprintf("mpd - [%s]", status[1]) //status[1] should now be stopped or paused
		return
	}

	message = "currentsong\n"
	mpdConn.Write([]byte(message))
	mpdConn.Read(reply)

	// Close connection to mpd
	mpdConn.Write([]byte("close\n"))

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
		np = artist + " -" + title
		np = strings.Replace(np, "&", "+", -1)
		return
	} else { //This is a nonfatal error.
		np = "Playlist is empty."
		return
	}
}
