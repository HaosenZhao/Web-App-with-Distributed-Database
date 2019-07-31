package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
)

//
//The text file is just to initialized the main page as a default blog, all other blogs will be in memory
//Page class is to save the blog, with title, body, and first line of the blog in the main page
type Page struct {
	Title string
	Body  string
	FirstLine string
}

//Pages save all the pages created in the blog web
type Pages struct{
	Pages []*Page
}

type PLock struct{
	sync.Mutex
}

//page_lst is a global variable of Pages
var page_lst Pages
var plock PLock

//initialize function initialize the web main page and insert the default blog, which is saved in a text file
//and all other data saved in local
func initialize(){
	page_lst.Pages = []*Page{}
	initial_page := Page{}
	initial_page.Title = "A Gallry Called Hope"
	path,_ := os.Getwd()
	body,err := ioutil.ReadFile(path+"/A Gallry Called Hope.txt")
	if err != nil{
		fmt.Println(err.Error())
	}
	body_str := string(body)
	initial_page.Body = body_str
	first_string := strings.SplitAfter(body_str, ".")[0] + "..."
	initial_page.FirstLine = first_string
	page_lst.Pages = append(page_lst.Pages, &initial_page)
}

//receive information through a decoder
func receive(decoder *gob.Decoder) {
	plock.Lock()
	page_lst = Pages{}
	err := decoder.Decode(&page_lst)
	if err != nil{
		fmt.Println(err.Error())
	}
	plock.Unlock()
}

//send information through a encoder
func send(encoder *gob.Encoder) {
	plock.Lock()
	err := encoder.Encode(page_lst)
	if err != nil{
		fmt.Println(err.Error())
	}
	plock.Unlock()
}


// --backend, which indicate the endpoint used to communicate with the backend.
//An endpoint consists of a hostname and port, separated by a colon;
//if the hostname is omitted, it is assumed to refer to the local host.
func main() {
	portPtr := flag.String("listen", "8090", "a string")
	flag.Parse()
	initialize()
	listener, error := net.Listen("tcp", ":"+*portPtr)
	if error != nil {
		fmt.Println(error.Error())
	}
	fmt.Println("Listen to",*portPtr,", wait for connection.")

	for{
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		fmt.Println("A new connection from",connection.LocalAddr())

		go func() {
			defer connection.Close()
			encoder := gob.NewEncoder(connection)
			decoder := gob.NewDecoder(connection)

			for {
				request, _ := bufio.NewReader(connection).ReadString('R')
				request = string(request)
				if request == "readR" {
					send(encoder)
				}
				if request == "writeR" {
					receive(decoder)
				}
				if request == "HiR"{
					fmt.Fprintf(connection, "ImAliveR")
				}
			}
		}()
	}
}