package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	// my packages
	"main/server"
)

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
			"detail": "Authentication credentials are invalid or missing.",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(notAuthReply)
}

func handler_get_session_fromid(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received GET request for session ID")
	// confirm is logged in

	path := strings.TrimPrefix(r.URL.Path, "/v1/sessions/")
	path = strings.TrimSuffix(path, "/")
	{
		if path == "" {
			http.Error(w, "Session ID is required", http.StatusBadRequest)
			return
		}

		sessionId := path
		fmt.Printf("Received request for session ID: %s\n", sessionId)
		w.Header().Set("Content-Type", "application/json")
		response := server.BaseReply{
			Status: "success",
			Data: server.SessionData{
				Username:    "testuser",
				Client_type: "direct",
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func handler_post_authenticate(w http.ResponseWriter, r *http.Request) { // return session_id on valid login
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Received POST request for authentication")

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

	if server.DefaultServer.AuthenticateLogin(data) {
		// Simulate successful authentication
		w.WriteHeader(http.StatusOK)
		response := server.BaseReply{
			Status: "success",
			Data:   map[string]string{"session_id": server.DefaultServer.GenerateClientSessionID(data.Username)},
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
	// This will match every request as a fallback handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if server.DefaultServer.IsRateLimited(r.RemoteAddr) {
			fmt.Println("Client is rate limited: ", r.RemoteAddr)
			http_too_many_requests(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/v1/sessions/") {
			if r.Method == http.MethodPost { // post request, user is trying to authenticate
				fmt.Println("Received POST request for session authentication")
				handler_post_authenticate(w, r)
				return
			}

			// client has to be authenticated beyond this point
			if r.Header.Get("Session") != server.DefaultServer.GenerateClientSessionID("testuser") {
				http_autherror_sessionid_invalid(w, r)
				fmt.Println("Invalid session ID provided")
				return // early return if session ID is invalid
			}

			// if we're here, client has a valid session header
			fmt.Println("Valid session ID provided, processing GET request")

			if r.Method == http.MethodGet {
				handler_get_session_fromid(w, r)
			}
			return
		} else {
			http_not_found(w, r)
			return
		}
	})

	fmt.Printf("Server running @ %s\n", server.DefaultServer.URL)
	http.ListenAndServe(":8080", nil)
}

// /v1/sessions/<session_id>/
