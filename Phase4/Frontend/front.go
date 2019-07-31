package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Frontend just handle different handler HTTP requests, no data stored
//The web app has 6 handler, edit, save, view, delete, create and the main page
//Some Handler is from golang official website
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

type Backend struct{
	serverPtr *string
	server net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
	status bool
	id int
}

//Encoder, decoder and server are used as a method to encode or decode object, no information stored
//var encoder *gob.Encoder
//var decoder *gob.Decoder
//var server net.Conn
//var serverPtr *string
var portPtr *string
var backends []Backend

//find_page is to return the pointer and index of page using a title
func find_page(title string, page_lst Pages)(*Page, int){
	for i := 0; i < len(page_lst.Pages); i++{
		if page_lst.Pages[i].Title == title{
			return page_lst.Pages[i],i
		}
	}
	return nil,99999
}


//createHandler is to create a new blog in the main apge
func createHandler(w http.ResponseWriter, r *http.Request){
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		title := r.FormValue("title")
		pag, _ := find_page(title, page_lst)
		counter := 0
		ori_title := title
		for pag != nil {
			counter++
			title = ori_title + "(" + strconv.Itoa(counter) + ")"
			pag, _ = find_page(title, page_lst)
		}
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}

//editHandler is to edit the page created
func editHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		title := r.URL.Path[len("/edit/"):]
		p, _ := find_page(title, page_lst)
		if p == nil {
			p = &Page{Title: title, Body: ""}
			page_lst.Pages = append(page_lst.Pages, p)
			send(title, &page_lst)
		}
		path, _ := os.Getwd()
		t, _ := template.ParseFiles(path + "/edit.html")
		t.Execute(w, p)
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}

//saveHandler is to save the change in blog after edit
func saveHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		title := r.URL.Path[len("/save/"):]
		body := r.FormValue("body")
		p, _ := find_page(title, page_lst)
		p.Body = body
		p.FirstLine = strings.SplitAfter(body, ".")[0] + "..."
		send(title, &page_lst)
		http.Redirect(w, r, "/index", http.StatusFound)
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}


//viewHandler is to view the blog
func viewHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		title := r.URL.Path[len("/view/"):]
		p, _ := find_page(title, page_lst)
		if p == nil {
			http.Redirect(w, r, "/index-delete", http.StatusFound)
			return
		}
		path, _ := os.Getwd()
		t, _ := template.ParseFiles(path + "/view.html")
		t.Execute(w, p)
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}

//deleteHandler is to delete the blog in the edit page
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		title := r.URL.Path[len("/delete/"):]
		_, index := find_page(title, page_lst)
		if index == 99999 {
			http.Redirect(w, r, "/index", http.StatusFound)
			return
		}
		page_lst.Pages = append(page_lst.Pages[:index], page_lst.Pages[index+1:]...)
		send(title, &page_lst)
		http.Redirect(w, r, "/index", http.StatusFound)
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}

//frontPageHandler is to create the index page
func frontPageHandler(w http.ResponseWriter, r *http.Request){
	var page_lst Pages
	if testAliveQ() {
		receive(&page_lst)
		path, _ := os.Getwd()
		t, _ := template.ParseFiles(path + "/index.html")
		err := t.Execute(w, page_lst)
		if err != nil {
			fmt.Println(err.Error())
		}
	}else{
		http.Redirect(w, r, "/index-delete", http.StatusFound)
	}
}

//frontPageErrorHandler is to redirect if the page is deleted by other client
func frontPageErrorHandler(w http.ResponseWriter, r *http.Request){
	var page_lst Pages
	receive(&page_lst)
	path,_ := os.Getwd()
	t,_ := template.ParseFiles(path+"/index-error.html")
	err := t.Execute(w,page_lst)
	if err != nil {
		fmt.Println(err.Error())
	}
}

//receive information from all backends and store in page_lst
func receive(page_lst *Pages){
	for i:=0; i < len(backends); i++{
		if backends[i].status == true{
			fmt.Fprintf(backends[i].server, "readR")
			err := backends[i].decoder.Decode(page_lst)
			if err != nil{
				fmt.Println(err.Error())
			}
		}
	}
}

//test if the quorum of backend remains functioning
func testAliveQ() bool{
	quoNum:= len(backends)/2+1
	counter:=0
	for i:=0; i< len(backends);i++{
		if backends[i].status==true{
			counter++
		}
		if counter>=quoNum{
			return true
		}
	}
	return false
}


//hash the title
func hash(title string) int{
	result:=0
	for i:=0; i < len(title); i++{
		result+=int(title[i])
	}
	result=result%1024
	return result
}

//find the absolute value
func abs(num int) int{
	if num <0{
		return -1*num
	}
	return num
}

//send information as page_lst to the backend with id closet to the hash number
func send(title string,page_lst *Pages) {
	hashInt := hash(title)
	var server net.Conn
	offset:=9999
	for i:=0;i < len(backends);i++{
		if backends[i].status == true{
			diff := abs(backends[i].id-hashInt)
			if diff < offset{
				server = backends[i].server
			}
			fmt.Fprintf(server, "writeR")
			time.Sleep(time.Millisecond*10)
			_ = backends[i].encoder.Encode(page_lst)
		}
	}
}

//testAlive the backends and reconnect if possible
func testAlive(index int){
	for {
		time.Sleep(time.Second)
		if backends[index].status==false{
			loc, _ := time.LoadLocation("UTC")
			now := time.Now().In(loc)
			fmt.Println("Detected failure on backend with id " + strconv.Itoa(backends[index].id) +" at " + now.String())
			_,err := net.Dial("tcp", *backends[index].serverPtr)
			fmt.Fprintf(backends[index].server, "loseR")
			if err == nil{
				backends[index].server,_ = net.Dial("tcp", *backends[index].serverPtr)
				backends[index].encoder = gob.NewEncoder(backends[index].server)
				backends[index].decoder = gob.NewDecoder(backends[index].server)
				backends[index].status = true
			}
		}else {
			backends[index].server.SetReadDeadline(time.Now().Add(time.Second * 3))
			fmt.Fprintf(backends[index].server, "HiR")
			time.Sleep(time.Millisecond*20)
			_, err := bufio.NewReader(backends[index].server).ReadString('R')
			if err != nil {
				backends[index].status=false
			}
		}
	}
}


//2 command line arguments provided
//listen indicates the port number to accept HTTP connections on, if unspecified, this will be 8080
//backend indicates the endpoints used to communicate with the backends, if no hostname provided, the hostname
//will be localhost, and default port number for backend is 8090.
func main() {
	portPtr = flag.String("listen", "8080", "a string")
	serverPtrStr := flag.String("backend", ":8090", "a string")
	flag.Parse()
	serverLst := strings.Split(*serverPtrStr, ",")
	backends = []Backend{}
	for i:=0; i < len(serverLst); i++ {
		var err error
		backends = append(backends, Backend{})
		backends[i].serverPtr= &serverLst[i]
		backends[i].server,err =net.Dial("tcp", *backends[i].serverPtr)
		if err == nil {
			backends[i].encoder = gob.NewEncoder(backends[i].server)
			backends[i].decoder = gob.NewDecoder(backends[i].server)
			backends[i].status = true
			fmt.Fprintf(backends[i].server, "idR")
			id, _ := bufio.NewReader(backends[i].server).ReadString('R')
			if id=="9999R"{
				backends[i].id = rand.Intn(1025)
				id = strconv.Itoa(backends[i].id)
				fmt.Fprintf(backends[i].server, id+"R")
			}else{
				backends[i].id, _= strconv.Atoi(id)
			}
		}else{
			backends[i].status = false
		}
		go testAlive(i)
	}
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/index", frontPageHandler)
	http.HandleFunc("/index-delete", frontPageErrorHandler)
	fmt.Println("Go to 127.0.0.1:"+*portPtr+"/index")
	log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
}