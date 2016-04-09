package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmcvetta/napping"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	sessn          napping.Session
	headers        http.Header
	tsport         http.Transport
	clnt           http.Client
	allowedMethods map[string]bool
)

type SMSMessage struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Message     string `json:"message"`
	Campaign    string `json:"campaign"`
	SharedPool  string `json:"sharedPool"`
}

type SMSResponse struct {
	Errors []json.RawMessage
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

	// REST connection setup
	tsport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	clnt = http.Client{Transport: tsport}

	allowedMethods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"PATCH":   true,
		"DELETE":  true,
		"OPTIONS": true,
	}

}

func GetAuthHeader(method string, uri string) (error, string) {

	// Authorization: MAC id="your API key", ts="1325376000", nonce="random-string", mac="base64-encoded-hash"

	if !allowedMethods[method] {
		msg := fmt.Sprintf("invalid method: %v", method)
		err := errors.New(msg)
		return err, ""
	}

	now := time.Now().Unix()
	tstr := strconv.FormatInt(now, 10)
	nonce := createNonce(now)
	var buffer bytes.Buffer
	buffer.WriteString("MAC id=")
	buffer.WriteString(fmt.Sprintf("%+q", apikey))
	buffer.WriteString(", ts=")
	buffer.WriteString(fmt.Sprintf("%+q", tstr))
	buffer.WriteString(", nonce=")
	buffer.WriteString(fmt.Sprintf("%+q", nonce))
	buffer.WriteString(", mac=")
	macstr := tstr + "\n"
	macstr += nonce + "\n"
	macstr += method + "\n" + uri + "\napi.smsglobal.com\n443\n\n"
	mac := computeHmac256(macstr, secret)
	buffer.WriteString(fmt.Sprintf("%+q", mac))

	return nil, buffer.String()

}

func createNonce(now int64) string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(now, 10))
	io.WriteString(h, strconv.FormatInt(rand.Int63(), 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func PostMsg(url string, payload interface{}, result interface{}) (error, *napping.Response) {

	e := httperr{}
	var (
		err  error
		resp *napping.Response
	)

	err, auth := GetAuthHeader("POST", "/v1/sms/")
	if err != nil {
		return err, nil
	}
	if debug {
		fmt.Println("auth: ", auth)
	}
	headers.Add("Authorization", auth)
	headers.Add("Accept", "application/json")

	sessn = napping.Session{
		Client: &clnt,
		Log:    debug,
		Header: &headers,
	}

	//
	// Send request to server
	//

	resp, err = sessn.Post(url, &payload, &result, &e)

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
