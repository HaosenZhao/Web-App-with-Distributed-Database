package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strconv"
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

type Backend struct{
	serverPtr *string
	id int
	pgs Pages
}

type BackLink struct {
	prev Backend
	next Backend
}

//page_lst is a global variable of Pages
var page_lst Pages
var plock PLock
var id int
var bklink BackLink
var bkends []Backend

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
	bklink= BackLink{prev:Backend{},next:Backend{}}
	bkends= []Backend{}
	first_string := strings.SplitAfter(body_str, ".")[0] + "..."
	initial_page.FirstLine = first_string
	page_lst.Pages = append(page_lst.Pages, &initial_page)
	id = 9999
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

//create the link list of backend
func createLk(lkArray []string){
	ip_id := map[int]string{}
	for i:=0; i< len(lkArray);i++{
		for j:=0; j < len(bkends); j++{
			if lkArray[i]==*bkends[j].serverPtr{
				ip_id[bkends[j].id]=lkArray[i]
			}
		}
	}
	var keys[]int
	for k := range ip_id{
		keys = append(keys, k)
	}
	sort.Ints(keys)
	index := 0
	for i:=0; i< len(keys); i++{
		if keys[i] == id{
			index=i
		}
	}
	ip := []string{}
	for _,k := range keys{
		ip=append(ip, ip_id[k])
	}
	if len(ip) == 0{
		return
	}
	if index !=0 && index != len(keys){
		bklink.prev.id = keys[index-1]
		bklink.prev.serverPtr = &ip[index-1]
		bklink.next.id = keys[index+1]
		bklink.next.serverPtr = &ip[index+1]
	}else if index == 0{
		bklink.prev.id = keys[len(keys)]
		bklink.prev.serverPtr = &ip[len(ip)]
		bklink.next.id = keys[index+1]
		bklink.next.serverPtr = &ip[index+1]
	}else{
		bklink.prev.id = keys[index-1]
		bklink.prev.serverPtr = &ip[index-1]
		bklink.next.id = keys[0]
		bklink.next.serverPtr = &ip[0]
	}
}

//transfer the data of previous backend if it is lost
func transfer(){
	for i := 0; i < len(bklink.prev.pgs.Pages); i++{
		page_lst.Pages=append(page_lst.Pages, bklink.prev.pgs.Pages[i])
	}
	return
}

// --listen, which indicate the endpoint used to communicate with the backend.
//An endpoint consists of a hostname and port, separated by a colon;
//if the hostname is omitted, it is assumed to refer to the local host.
// --backend, which indicate the endpoint of other backends seperated by comma
func main() {
	portPtr := flag.String("listen", "8090", "a string")
	serverPtrStr := flag.String("backend", ":8090", "a string")
	flag.Parse()
	serverLst := strings.Split(*serverPtrStr, ",")
	initialize()
	listener, error := net.Listen("tcp", ":"+*portPtr)
	if error != nil {
		fmt.Println(error.Error())
	}
	fmt.Println("Listen to",*portPtr,", wait for connection.")
	createLk(serverLst)
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
				if request == "readR" {
					send(encoder)
				} else if request == "writeR" {
					receive(decoder)
				} else if request == "HiR"{
					fmt.Fprintf(connection, "ImAliveR")
				} else if request == "idR"{
					fmt.Fprintf(connection, strconv.Itoa(id)+"R")
				} else if request == "loseR"{
					transfer()
				}else{
					if len(request)>1{
						request = request[:len(request)-1]
						id, _ = strconv.Atoi(request)
						fmt.Println("ID: "+strconv.Itoa(id))
					}
				}
			}
		}()
	}
}