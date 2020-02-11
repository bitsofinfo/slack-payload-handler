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
	listenPort int // port to listen on (flag opt)
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
		//Methods("POST").
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

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	// first lets get the credentials off the request
	validated, err := validateRequestSignature(req)
	if err != nil || !validated {
		writeHTTPResponse(resWriter, "Bad Request: security check failed", http.StatusBadRequest)
		return
	}

	var x map[string]interface{}
	err = json.Unmarshal([]byte(req.FormValue("payload")), &x)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	var actionValuesArr []string

	var actionsArr []interface{} = x["actions"].([]interface{})
	//fmt.Printf("%v", actionsArr[0]["value"])

	for i := 0; i < len(actionsArr); i++ {
		var actionMap map[string]interface{} = actionsArr[i].(map[string]interface{})
		fmt.Printf("%v\n", actionMap["value"])
		var theval string = actionMap["value"].(string)
		actionValuesArr = append(actionValuesArr, theval)
	}

	x["action_values"] = actionValuesArr

	resWriter.Header().Set("Content-Type", "application/json")
	resWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(resWriter).Encode(x)

}
