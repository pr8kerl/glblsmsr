package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"time"
)

var (
	version string = "0.1."
	commit  string = "unset"
	cfgfile string = "config.json"
	message string
	number  string
	u       string
	debug   bool
)

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug")
	flag.StringVarP(&message, "message", "m", "", "the message to be sent")
	flag.StringVarP(&number, "number", "n", "", "the mobile phone number to send to")
	flag.StringVarP(&u, "url", "u", "http://localhost:8951", "url of the gosms service")

}

func main() {

	// commit is set by go build -ldflags in Makefile
	version = version + commit
	flag.Parse()
	u += "/api/sms/"

	if number == "" {
		fmt.Println("required option --number missing.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if message == "" {
		fmt.Println("required option --message missing.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	const tlayout = time.RFC3339
	t := time.Now()
	tstamp := t.Format(tlayout)

	if debug {
		fmt.Println("debug enabled.\n")
		fmt.Printf("url: %s\n", u)
		fmt.Printf("number: %s\n", number)
		message = tstamp + " " + message
		fmt.Printf("message: %s\n", message)
	}

	msg := SMSMessage{}
	msg.Mobile = number
	msg.Message = message
	res := SMSResponse{}
	pload := json.RawMessage{}
	pload, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error marshalling payload : %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("msg: %s\n", pload)
	err, resp := PostMsg(u, &pload, &res)
	if err != nil {
		fmt.Printf("err : %s\n", err)
		//		fmt.Printf("%s : %s\n", resp.HttpResponse().Status, resp.RawText())
		os.Exit(1)
	}
	//fmt.Printf("OK : %s\n", res.Message)
	fmt.Printf("%s : %s\n", resp.HttpResponse().Status, resp.RawText())

	os.Exit(0)
}
