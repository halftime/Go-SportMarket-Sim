package main

import "main/server"

func main() {
	go server.SetupWebSocketServer()
	server.SetupHttpServer()
}
