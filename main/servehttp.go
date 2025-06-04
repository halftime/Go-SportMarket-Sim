package main

// https://api.sportmarket.com/v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	// my packages

	"main/server"
)

func is_client_rate_limited(w http.ResponseWriter, r *http.Request) bool {
	// Check if the client is rate limited
	clientId := r.RemoteAddr //r.Header.Get("session_id")
	if server.DefaultServer.IsRateLimited(clientId) {
		return true
	}
	return false
}

func http_too_many_requests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API limit exceeded (429)", http.StatusTooManyRequests)
}

func http_not_found(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found (404)", http.StatusNotFound)
}

func http_autherror_sessionid_invalid(w http.ResponseWriter, r *http.Request) {
	notAuthReply := server.ErrorReply{
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

	var data server.LoginRequest
	missingParameters := make(map[string][]string)

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || r.Body == nil || (data.Username == "" || data.Password == "") {
		if data.Username == "" {
			missingParameters["username"] = []string{"This field is required"}
		}
		if data.Password == "" {
			missingParameters["password"] = []string{"This field is required"}
		}

		missingParamReply := server.ErrorReply{
			Status: "error",
			Code:   "validation_error",
			Data:   missingParameters,
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(missingParamReply)
		return
	}

	if data.Username == "testuser" && data.Password == "testpass" {
		w.WriteHeader(http.StatusOK)
		response := server.ReplyEnvelope{
			Status: "success",
			Data:   map[string]string{"session_id": "testsessionid"},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	failedAuthReply := server.ErrorReply{
		Status: "error",
		Code:   "authentication_failed",
		Data:   "Authentication failed",
	}
	w.WriteHeader(http.StatusUnauthorized) // (401)
	json.NewEncoder(w).Encode(failedAuthReply)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if is_client_rate_limited(w, r) {
			// If the client is rate limited, return a 429 Too Many Requests response
			fmt.Println("Client is rate limited:", r.RemoteAddr)
			http_too_many_requests(w, r)
			return
		}

		// Only continue if rate_limit_passthrough did not write a response
		switch r.URL.Path {
		case "/v1/sessions/":
			handler_authenticate(w, r)

		case "/v1/betslips/":
			http_autherror_sessionid_invalid(w, r)

		default:
			http_not_found(w, r)
		}
	})

	fmt.Printf("Server running @ %s\n", server.DefaultServer.URL)
	http.ListenAndServe(":8080", nil)
}
