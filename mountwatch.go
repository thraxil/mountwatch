package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

type mount struct {
	Path   string `json:"path"`
	Metric string `json:"metric"`
}

type config struct {
	Mounts []mount `json:"mounts"`
}

var (
	configfile      = flag.String("config", "config.json", "path to config JSON file")
	graphiteAddress = flag.String("graphite", "", "Graphite service address (example: 'localhost:2003')")
	interval        = flag.Int64("interval", 60, "Check interval")
	prefix          = flag.String("prefix", "stats", "Prefix")
)

func check(prefix string, mounts []mount) *bytes.Buffer {
	now := int32(time.Now().Unix())
	buffer := bytes.NewBufferString("")
	for _, m := range mounts {
		_, err := ioutil.ReadDir(m.Path)
		if err != nil {
			fmt.Fprintf(buffer, "%s.%s %d %d\n", prefix, m.Metric, 1, now)
		} else {
			fmt.Fprintf(buffer, "%s.%s %d %d\n", prefix, m.Metric, 0, now)
		}
	}
	return buffer
}

func submit(address string, buffer *bytes.Buffer) {
	var clientGraphite net.Conn
	if address != "" {
		var err error
		clientGraphite, err = net.Dial("tcp", address)
		if clientGraphite != nil {
			// Run this when we're all done, only if clientGraphite was opened.
			defer clientGraphite.Close()
		}
		if err != nil {
			log.Printf(err.Error())
		}
	}
	clientGraphite.Write(buffer.Bytes())
}

func monitor(prefix string, mounts []mount, address string, interval int64) {
	for {
		buffer := check(prefix, mounts)
		fmt.Print(buffer)
		submit(address, buffer)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {
	flag.Parse()
	file, err := ioutil.ReadFile(*configfile)
	if err != nil {
		log.Fatal(err)
	}
	c := config{}
	err = json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal(err)
	}
	monitor(*prefix, c.Mounts, *graphiteAddress, *interval)
}
