package main
import (
	"fmt"
	"encoding/gob"
	"net"
	"log"
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
	fmt.Println("Client")
	//create structure object

	studentEncode := Pages{Pages: []*Page{}}
	studentEncode.Pages = append(studentEncode.Pages, &Page{"life of raccoon", "this is life of raccoon", "asd"})

	fmt.Println("start client");
	// dial TCP connection
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	//Create encoder object, We are passing connection object in Encoder
	encoder := gob.NewEncoder(conn)
	// Encode Structure, IT will pass student object over TCP connection
	encoder.Encode(studentEncode)

	dec := gob.NewDecoder(conn)
	// create blank student object
	p := &Pages{}
	// decode serialize data
	dec.Decode(p)

	fmt.Println(p.Pages[0].FirstLine)

	// close connection
	conn.Close()
	fmt.Println("done");
}