package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

/*
Example of what we're doing here:
We send a post request to e.g: http://leia-l:9000/jsonrpc.js

The body looks like e.g. (for pause):
{"id":1,"method":"slim.request","params":["00:04:20:2b:76:f6",["pause"]]}
e.g. (for play):
{"id":1,"method":"slim.request","params":["00:04:20:2b:76:f6",["play"]]}
*/

// lmsHostname is the hostname of the LMS server; here I am using a local DNS entry
// that is in my /etc/hosts file
const lmsHostname = "leia-l"

// default id of the player(?), used in LmsService struct
const defaultId = 1

// default method, used in LmsService struct
const defaultMethod = "slim.request"

// default volume change amount
const volumeUnit int = 5

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

func volumeString(positive bool) string {
	// pm is a plus or a minus
	// i.e. pm == "+" or pm == "-"
	var pm string

	if positive {
		pm = "+"
	} else {
		pm = "-"
	}

	return pm + strconv.Itoa(volumeUnit)
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
		Params: []any{playerMacAddress, []string{"mixer", "volume", volumeString(false)}},
	}
}

// LmsVolumeUp is a helper function to create a LmsService struct for a "volume up" operation
func LmsVolumeUp(playerMacAddress string) LmsService {
	return LmsService{
		Id:     defaultId,
		Method: defaultMethod,
		Params: []any{playerMacAddress, []string{"mixer", "volume", volumeString(true)}},
	}
}

// LmsPost is a helper function to send a POST request to the LMS server
// with the given LmsService struct
// returning the HTTP status code and body text
func LmsPost(playerMacAddress string, fn func(string) LmsService) (int, string) {
	// where to POST to
	postUrl := "http://" + lmsHostname + ":9000/jsonrpc.js"
	// create the JSON payload
	jsonData, err := json.Marshal(fn(playerMacAddress))
	// bork on error
	checkErr(err)

	// create POST request
	postRequest, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	// bork on error
	checkErr(err)
	// set the content type header to JSON
	postRequest.Header.Set("Content-Type", "application/json")

	// client struct
	client := &http.Client{}
	// send HTTP request
	httpResponse, err := client.Do(postRequest)
	// bork on error
	checkErr(err)
	// close when done
	defer httpResponse.Body.Close()

	// read the body text
	bodyText, err := io.ReadAll(httpResponse.Body)
	// bork on error
	checkErr(err)
	// return the HTTP status code and body text
	return httpResponse.StatusCode, string(bodyText)
}
