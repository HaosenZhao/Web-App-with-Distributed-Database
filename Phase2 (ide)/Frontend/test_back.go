package main

import (
	"fmt"
	"net"
	"encoding/gob"
)

type Page struct {
	Title string
	Body  string
	FirstLine string
}

//Pages save all the pages created in the blog web
type Pages struct{
	Pages []*Page
}


func main() {
	fmt.Println("Server")
	// start TCP server and listen on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
		panic(err)
	}
	conn, err := ln.Accept()
	dec := gob.NewDecoder(conn)
	// create blank student object
	p := &Pages{}
	// decode serialize data
	dec.Decode(p)
	// print
	encoder := gob.NewEncoder(conn)
	// Encode Structure, IT will pass student object over TCP connection
	encoder.Encode(*p)
	// close connection for that client
	conn.Close()

}