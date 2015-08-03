package main

import (
	"encoding/json"
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
	debug   bool
	server  string = "api.smsglobal.com"
	uri     string = "/v1/sms/"
	apikey  string = "b100085ef1a79e5c29e913c9084e0e89"
	secret  string = "0e0fe8de088bced022726f18cb01e6c5"
	from    string
)

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug")
	flag.StringVarP(&message, "message", "m", "", "the message to be sent")
	flag.StringVarP(&number, "number", "n", "", "the mobile phone number to send to")
	flag.StringVarP(&from, "from", "f", "smsglobal", "string identifier of who is sending the msg")

}

func main() {

	// commit is set by go build -ldflags in Makefile
	version = version + commit
	flag.Parse()
	u = "https://" + server + uri

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

	if debug {
		fmt.Println("debug enabled.\n")
		fmt.Printf("url: %s\n", u)
		fmt.Printf("number: %s\n", number)
		fmt.Printf("message: %s\n", message)
	}

	msg := SMSMessage{}
	msg.Origin = from
	msg.Destination = number
	msg.Message = message
	res := json.RawMessage{}
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
	printResponse(res)
	fmt.Printf("%s : %s\n", resp.HttpResponse().Status, resp.RawText())

	os.Exit(0)
}
