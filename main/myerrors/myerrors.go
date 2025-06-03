package myerrors

type ReplyEnvelope struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type ErrorReply struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Data   any    `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Api_not_auth_reply = ErrorReply{
	Status: "error",
	Code:   "auth_error",
	Data: map[string]string{
		"detail": "Authentication credentials were not provided.",
	},
}

// example for non authenticated request
// {"status":"error","code":"auth_error","data":{"detail":"Authentication credentials were not provided."}}

// example for failed auth attempt
// {"status":"error","code":"authentication_failed","data":"Authentication failed"}
