package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

var (
	version string = "0.1."
	commit  string = "unset"
	cfgfile string = "config.json"
	message string
	number  string
	u       string
)

func init() {

	flag.StringVarP(&message, "message", "m", "hello", "the message to be sent")
	flag.StringVarP(&number, "number", "n", "", "the mobile phone number to send to")
	flag.StringVarP(&u, "url", "u", "http://localhost:8951", "url of the gosms service")

}

func main() {

	// commit is set by go build -ldflags in Makefile
	version = version + commit
	flag.Parse()
	u += "/api/sms"
	fmt.Printf("url: %s\n", u)
	fmt.Printf("number: %s\n", number)
	fmt.Printf("message: %s\n", message)

	msg := SMSMessage{number, message}
	res := SMSResponse{}

	err, _ := PostMsg(u, msg, &res)
	if err != nil {
		fmt.Printf("ERR : %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("OK : %s\n", res.Message)

	os.Exit(0)
}
