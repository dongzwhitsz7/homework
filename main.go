package main

import "dongzw/dongzwhom/module2"

func main() {
	server := module2.NewHttpServer()
	server.Serve()
}
