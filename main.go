package main

import (
	"dongzw/dongzwhom/http_server"
)

func main() {
	server := http_server.NewHttpServer()
	server.Serve()
}
