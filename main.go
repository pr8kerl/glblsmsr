package main

import (
	"bufio"
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
)

var (
	message string
	number  string
	u       string
	debug   bool
	apikey  string
	secret  string
	from    string
)

const (
	server   string = "api.smsglobal.com"
	uri      string = "/v1/sms/"
	glblkey  string = "GLBLKEY"
	glblscrt string = "GLBLSCRT"
)

func init() {

	flag.BoolVarP(&debug, "debug", "d", false, "enable debug")
	flag.StringVarP(&message, "message", "m", "", "the message to be sent")
	flag.StringVarP(&number, "number", "n", "", "the mobile phone number to send to")
	flag.StringVarP(&from, "from", "f", "smsglobal", "string identifier of who is sending the msg")

}

func main() {

	flag.Parse()

	if number == "" {
		fmt.Printf("required option --number missing.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if message == "" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message += scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error reading stdin:", err)
		}
	}
	if message == "" {
		fmt.Printf("no message provided.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	apikey = os.Getenv(glblkey)
	if apikey == "" {
		fmt.Printf("missing environment variable %s\n", glblkey)
		os.Exit(1)
	}
	secret = os.Getenv(glblscrt)
	if secret == "" {
		fmt.Printf("missing environment variable %s\n", glblscrt)
		os.Exit(1)
	}

	u = "https://" + server + uri

	if debug {
		fmt.Printf("debug enabled.\n")
		fmt.Printf("url: %s\n", u)
		fmt.Printf("number: %s\n", number)
		fmt.Printf("message: %s\n", message)
	}

	requ := SMSMessage{}
	resu := SMSResponse{}
	requ.Origin = from
	requ.Destination = number
	requ.Message = message

	err, resp := PostMsg(u, &requ, &resu)
	if err != nil {
		fmt.Printf("%s : %s\n", resp.HttpResponse().Status, err)
		os.Exit(1)
	}

	if resp.Status() != 201 {
		fmt.Printf("ERR %s\n", resp.HttpResponse().Status)
		printResponse(&resu)
		os.Exit(resp.Status())
	} else {
		head := resp.HttpResponse().Header
		msgid := head.Get("Location")
		fmt.Printf("%s %s\n", resp.HttpResponse().Status, msgid)
	}
	os.Exit(0)
}
