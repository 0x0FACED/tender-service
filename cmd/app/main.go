package main

import "github.com/0x0FACED/tender-service/internal/app/server"

func main() {
	if err := server.Start(); err != nil {
		panic("Server didnt start" + err.Error())
	}
}
