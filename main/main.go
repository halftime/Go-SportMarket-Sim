package main

// https://api.sportmarket.com/v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	// my packages
	"main/myerrors"
)

func return_http_too_many_requests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API limit exceeded (429)", http.StatusTooManyRequests)
}

func return_http_not_found(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found (404)", http.StatusNotFound)
}

func return_http_not_authorized(w http.ResponseWriter, r *http.Request) {
	notAuthReply := myerrors.ErrorReply{
		Status: "error",
		Code:   "auth_error",
		Data: map[string]string{
			"detail": "Authentication credentials were not provided.",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(notAuthReply)
}

func handler_get(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "This is a GET request handler. Use POST to submit data.")
}

func handler_authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data myerrors.LoginRequest
	missingParameters := make(map[string][]string)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || r.Body == nil || (data.Username == "" || data.Password == "") {
		if data.Username == "" {
			missingParameters["username"] = []string{"This field is required"}
		}
		if data.Password == "" {
			missingParameters["password"] = []string{"This field is required"}
		}

		notAuthReply := myerrors.ErrorReply{
			Status: "error",
			Code:   "validation_error",
			Data:   missingParameters,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(notAuthReply)
		return
	}

	if data.Username == "testuser" && data.Password == "testpass" {
		w.WriteHeader(http.StatusOK)
		response := myerrors.ReplyEnvelope{
			Status: "success",
			Data:   map[string]string{"session_id": "testsessionid"},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	failedAuthReply := myerrors.ErrorReply{
		Status: "error",
		Code:   "authentication_failed",
		Data:   "Authentication failed",
	}
	w.WriteHeader(http.StatusUnauthorized) // (401)
	json.NewEncoder(w).Encode(failedAuthReply)
}

func main() {
	http.HandleFunc("/v1/sessions/", handler_authenticate) // Handle /v1/sessions/ with 429 Too Many Requests

	http.HandleFunc("/v1", return_http_not_found)                  // Handle /v1 endpoint with 404 Not Found
	http.HandleFunc("/v1/betslips", return_http_too_many_requests) // Handle /v1/betslips

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
