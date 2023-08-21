package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/toke/golang-callmon/fritzbox"
)

func handleMessages(msgchan <-chan fritzbox.FbEvent) {
	for {
		ev := <-msgchan

		jsonm, _ := json.Marshal(&ev)
		fmt.Printf("Some JSON: %s\n", jsonm)

		if ev.EventName == fritzbox.CALL {
			fmt.Printf("%s Event: %s->%s\n", ev.EventName, ev.Source, ev.Destination)
		} else if ev.EventName == fritzbox.RING {
			fmt.Printf("%s Event: %s->%s\n", ev.EventName, ev.Source, ev.Destination)
		} else {
			fmt.Printf("! %s\n", ev)
		}
	}
}

func mainloop(host string) {

	recv := make(chan fritzbox.FbEvent)

	c, err := new(fritzbox.CallmonHandler).Connect(host, recv)
	if err != nil {
		return
	}

	defer c.Close()

	if c.Connected {
		go handleMessages(recv)

		// Inject a test message
		f := c.Parse("06.08.14 14:52:26;CALL;1;10;50000001;012344567;SIP1;\r\n")
		recv <- f

		c.Loop()
	}
}

func main() {
	arg := os.Args

	host := "fritz.box"
	if len(arg) > 1 && arg[1] != "" {
		host = arg[1]
	}

	for {
		mainloop(host)
		time.Sleep(1 * time.Second)
		fmt.Println("reconnect...")
	}
}
