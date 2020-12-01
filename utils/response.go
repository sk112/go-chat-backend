package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	host = "https://desolate-brook-30710.herokuapp.com"
)

//ResponseModel to send responses
type ResponseModel struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Send : single method using which all requests are to be sent
func Send(w http.ResponseWriter, log string, code int, message string, data interface{}) {
	fmt.Println(time.Now().Local().Format("2006.01.02 15:04:05"), log)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", host)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(code)

	// create a json
	response := ResponseModel{
		Message: message,
		Data:    data,
	}

	marshalledJSON, err := json.Marshal(response)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(marshalledJSON)
}

// SendPreflightResponse : for handling CORS policy -- for debugging with react hot reload
// NOTE a proxy can be set up in react app to redirect the calls to the port on golang which server is running
func SendPreflightResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println(time.Now().Local().Format("2006.01.02 15:04:05"), "response: sending preflight response for CORS")
	w.Header().Add("Access-Control-Allow-Origin", host)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, Authorization")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	w.WriteHeader(http.StatusOK)
}
