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
	sessn   napping.Session
	headers http.Header
	tsport  http.Transport
	clnt    http.Client
)

type SMSMessage struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Message     string `json:"message"`
	Campaign    string `json:"campaign"`
	SharedPool  string `json:"sharedPool"`
}

/*
origin	Where the SMS appears to come from. 4-11 characters A-Za-z0-9 if alphanumeric; 4-15 digits if numeric (if set, set sharedPool to null)	string
destination	Destination mobile number. 4-15 digits	string
message	The SMS message. If longer than 160 characters (GSM) or 70 characters (Unicode), splits into multiple SMS	string
campaign	The campaign the message is associated with (optional)	related
sharedPool	The shared pool to use (if set, set origin to null)	related
*/

type httperr struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

// send: 'GET /v1/balance HTTP/1.1\r\nHost: api.smsglobal.com\r\nAccept-Encoding: identity\r\nAccept: application/json\r\nAuthorization: MAC id="27a657ff3aec742ddca08e3d918f9ccd",ts="1372730113",nonce="bacd7eabff1864491e2071b2d3d6ae5f",mac="U8atdR0odGcTtmr2u7yaSvR7C1L1Qf3LjIZlqRn5MlM="\r\nUser-Agent: SMS Python Client\r\n\r\n'

/*
Authorization header fields
The value of the header is made up of the following components.
id

Your API key, issued to you by SMSGlobal. You can create an API key in your MXT account.
ts

The Unix timestamp of the time you made the request. We allow a slight buffer on this in case of any time sync issues.
nonce

A randomly generated string of your choice. Ensure it is unique to each request, and no more than 32 characters long.
To prevent replay attacks, a single nonce can only be used once.
mac

This is the base 64 encoded hash of the request.
Calculating the mac hash

The hash is a SHA-256 digest of a concatenation of a series of strings related to the request. The string is:
Timestamp
Nonce
HTTP request method
HTTP request URI
HTTP host
HTTP port
Optional extra data

Example string for POST to /v1/sms endpoint:
1325376000\n
random-string\n
POST\n
/v1/sms/\n
api.smsglobal.com\n
443\n
\n



*/

func init() {

	headers = make(http.Header)

	// REST connection setup
	tsport = http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	clnt = http.Client{Transport: &tsport}

}

func GetAuthHeader() string {

	// Authorization: MAC id="your API key", ts="1325376000", nonce="random-string", mac="base64-encoded-hash"

	now := time.Now().Unix()
	tstr := strconv.FormatInt(now, 10)
	//tstr := fmt.Sprintf("%f", tstamp)
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
	macstr += "POST\n/v1/sms/\napi.smsglobal.com\n443\n\n"
	mac := computeHmac256(macstr, secret)
	//	mac := hmac.New(sha256.New, []byte(secret))
	//	mac.Write([]byte(macstr))
	buffer.WriteString(fmt.Sprintf("%+q", mac))
	return buffer.String()

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

func PostMsg(u string, pload interface{}, res interface{}) (error, *napping.Response) {

	auth := GetAuthHeader()
	fmt.Println("auth: ", auth)
	headers.Add("Authorization", auth)

	sessn = napping.Session{
		Client: &clnt,
		Log:    debug,
		Header: &headers,
	}

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
