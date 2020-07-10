package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var (
	listenPort    int          // port to listen on (flag opt)
	debugRequest  bool = false // dumps raw request from slack to STDOUT
	debugResponse bool = false // dumps returned/modified response to STDOUT
)

// Struct for JSON we retun to caller
type actionsResponse struct {
	Whatever string `json:"whatever"`
}

type mutatedPayload struct {
	ActionsPressed []string
}

func init() {

	// cmd line args
	flag.IntVar(&listenPort, "listen-port", 8080, "Optional, port to listen on, default 8080")
	flag.BoolVar(&debugRequest, "debug-request", false, "Optional, print requests to STDOUT, default false")
	flag.BoolVar(&debugRequest, "debug-response", false, "Optional, print responses to STDOUT, default false")

	// logging options
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {

	flag.Parse()

	// setup our REST routes
	router := mux.NewRouter()
	router.Path("/").
		Methods("POST").
		Schemes("http").
		HandlerFunc(ProcessSlackRequest)

	// fire up the http server
	srv := &http.Server{
		Handler:      router,
		Addr:         (":" + strconv.Itoa(listenPort)),
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// writes an HTTP response w/ code + json result
func writeHTTPResponse(resWriter http.ResponseWriter, result string, httpStatus int) {
	resWriter.Header().Set("Content-Type", "application/json")
	resWriter.WriteHeader(httpStatus)
	json.NewEncoder(resWriter).Encode(&actionsResponse{Whatever: result})
}

// validates the request signature (TODO)
func validateRequestSignature(req *http.Request) (bool, error) {
	return true, nil
}

// ProcessSlackRequest ... http handler for processing the inbound slack POST payload
func ProcessSlackRequest(resWriter http.ResponseWriter, req *http.Request) {

	// Save a copy of this request for debugging?
	if debugRequest {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}

	// first lets get the credentials off the request
	validated, err := validateRequestSignature(req)
	if err != nil || !validated {
		writeHTTPResponse(resWriter, "Bad Request: security check failed", http.StatusBadRequest)
		return
	}

	// extract json from the payload
	var jsonPayload map[string]interface{}
	var rawData []byte = []byte(req.FormValue("payload"))
	if rawData != nil && len(rawData) > 0 {

		if debugRequest {
			fmt.Printf("\nJSON RECEIVED: \n%s\n\n", rawData)
		}

		err = json.Unmarshal(rawData, &jsonPayload)
		if err != nil {
			fmt.Printf("Could not parse Slack 'payload' into JSON: \n%v\n\n", err)
		}
	} else {
		log.Error("HTTP POST contained no FormData body where payload=[data]")
		return
	}

	if jsonPayload["actions"] == nil {
		log.Error("HTTP POST JSON contained no 'actions' element'")
		return
	}

	var actionValuesArr []string

	var actionsArr []interface{} = jsonPayload["actions"].([]interface{})

	for i := 0; i < len(actionsArr); i++ {
		var actionMap map[string]interface{} = actionsArr[i].(map[string]interface{})
		var theval string = actionMap["value"].(string)
		actionValuesArr = append(actionValuesArr, theval)
	}

	jsonPayload["action_values"] = actionValuesArr

	if debugResponse {
		var jsonStr, err = json.Marshal(jsonPayload)
		if err != nil {
			fmt.Printf("Could not parse Slack 'payload' into JSON: \n %v", err)
		} else if jsonStr != nil {
			fmt.Printf("RESPONSE: \n%s", jsonStr)
		}

	}

	log.Info(jsonPayload)

	resWriter.Header().Set("Content-Type", "application/json")
	resWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(resWriter).Encode(jsonPayload)

}
