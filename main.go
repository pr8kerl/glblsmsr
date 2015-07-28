package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

var (
	version  string = "0.1."
	commit   string = "unset"
	cfgfile  string = "config.json"
	message  string
	number   string
	hostname string = "localhost"
	port     int    = 8951
)

func init() {

	flag.StringVar(&message, "m", "message", "the message to be sent")
	flag.StringVar(&number, "n", "number", "the mobile phone number to send to")
	flag.StringVar(&hostname, "h", "hostname", "hostname of the gosms service")
	flag.IntVar(&port, "port", 8951, "port of the gosms service")

}

func main() {

	// commit is set by go build -ldflags in Makefile
	version = version + commit
	flag.Parse()
	fmt.Println("hostname: ", hostname)
	fmt.Println("port: ", port)
	fmt.Println("number: ", number)
	fmt.Println("message: ", message)

	os.Exit(0)
}
