package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

/*
Example of what we're doing here:
We send a post request to e.g: http://leia-l:9000/jsonrpc.js

The body looks like e.g. (for pause):
{"id":1,"method":"slim.request","params":["00:04:20:2b:76:f6",["pause"]]}
e.g. (for play):
{"id":1,"method":"slim.request","params":["00:04:20:2b:76:f6",["play"]]}
*/

// lmsHostnameConst is the hostname of the LMS server; here I am using a local DNS entry
// that is in my /etc/hosts file
const lmsHostnameConst = "leia-l"

const defaultId = 1
const defaultMethod = "slim.request"

// form of the JSON we send to the LMS server
type LmsService struct {
	Id     int    `json:"id"`
	Method string `json:"method"`
	Params []any  `json:"params"`
}

// checkErr will call log.Fatal if the error is not nil i.e. there is an error
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// LmsSimple is a helper function to create a LmsService struct for a simple operation
// with a given player MAC address and verb (e.g. "play" or "pause")
func LmsSimple(playerMacAddress string, verb string) LmsService {
	/*
		Looks like e.g:

		{
			"id": 1,
			"method": "slim.request",
			"params": [
				"00:04:20:2b:76:f6",
				[
					"pause"
				]
			]
		}
	*/

	return LmsService{
		Id:     defaultId,
		Method: defaultMethod,
		Params: []any{playerMacAddress, []string{verb}},
	}
}

// LmsPlay is a helper function to create a LmsService struct for a "play" operation
func LmsPlay(playerMacAddress string) LmsService {
	return LmsSimple(playerMacAddress, "play")
}

// LmsPause is a helper function to create a LmsService struct for a "pause" operation
func LmsPause(playerMacAddress string) LmsService {
	return LmsSimple(playerMacAddress, "pause")
}

// LmsVolumeDown is a helper function to create a LmsService struct for a "volume down" operation
func LmsVolumeDown(playerMacAddress string) LmsService {
	return LmsService{
		Id:     defaultId,
		Method: defaultMethod,
		Params: []any{playerMacAddress, []string{"mixer", "volume", "-10"}},
	}
}

// LmsVolumeUp is a helper function to create a LmsService struct for a "volume up" operation
func LmsVolumeUp(playerMacAddress string) LmsService {
	return LmsService{
		Id:     defaultId,
		Method: defaultMethod,
		Params: []any{playerMacAddress, []string{"mixer", "volume", "+10"}},
	}
}

// LmsPost is a helper function to send a POST request to the LMS server
// with the given LmsService struct
// returning the HTTP status code and body text
func LmsPost(playerMacAddress string, fn func(string) LmsService) (int, string) {
	postUrl := "http://" + lmsHostnameConst + ":9000/jsonrpc.js"
	jsondata, err := json.Marshal(fn(playerMacAddress))
	checkErr(err)

	r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsondata))
	checkErr(err)
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	checkErr(err)
	defer res.Body.Close()

	bodyText, err := io.ReadAll(res.Body)
	checkErr(err)
	return res.StatusCode, string(bodyText)
}
