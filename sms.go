package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
	"log"
	"net/http"
)

var (
	sessn   napping.Session
	headers http.Header
)

type SMSMessage struct {
	Mobile  string `json:"mobile"`
	Message string `json:"message"`
}

type SMSResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type httperr struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

func init() {

	headers = make(http.Header)

	sessn = napping.Session{
		Log:    debug,
		Header: &headers,
	}

}

func PostMsg(u string, pload interface{}, res interface{}) (error, *napping.Response) {

	//
	// Send request to server
	//
	e := httperr{}
	var (
		err  error
		resp *napping.Response
	)

	resp, err = sessn.Post(u, &pload, &res, &e)

	if err != nil {
		return err, resp
	}
	if resp.Status() >= 300 {
		return errors.New(e.Message), resp
	} else {
		// all is good in the world
		return nil, resp
	}
}

func printResponse(input interface{}) {

	jsonresp, err := json.MarshalIndent(&input, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonresp))

}
