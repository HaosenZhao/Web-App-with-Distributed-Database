package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"html/template"
	"os"
	"strconv"
	"strings"
	"time"
)

//The web app has 6 handler, edit, save, view, delete, create and the main page
//Some Handler is from
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

var encoder *gob.Encoder
var decoder *gob.Decoder
var server net.Conn
var portPtr *string

//find_page is to return the pointer and index of page using a title
func find_page(title string, page_lst Pages)(*Page, int){
	for i := 0; i < len(page_lst.Pages); i++{
		if page_lst.Pages[i].Title == title{
			return page_lst.Pages[i],i
		}
	}
	return nil,999
}


//createHandler is to create a new blog in the main apge
func createHandler(w http.ResponseWriter, r *http.Request){
	var page_lst Pages
	receive(&page_lst)
	title := r.FormValue("title")
	pag,_ := find_page(title,page_lst)
	counter:=0
	ori_title := title
	for pag != nil{
		counter++
		title = ori_title+"("+strconv.Itoa(counter)+")"
		pag,_ = find_page(title,page_lst)
	}
	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
}

//editHandler is to edit the page created
func editHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	receive(&page_lst)
	title := r.URL.Path[len("/edit/"):]
	p,_ := find_page(title,page_lst)
	if p == nil{
		p =  &Page{Title:title, Body:""}
		page_lst.Pages = append(page_lst.Pages, p)
		send(&page_lst)
	}
	path,_ := os.Getwd()
	t,_ := template.ParseFiles(path+"/Frontend/edit.html")
	t.Execute(w,p)
}

//saveHandler is to save the change in blog after edit
func saveHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	receive(&page_lst)
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p, _ := find_page(title,page_lst)
	p.Body = body
	p.FirstLine = strings.SplitAfter(body, ".")[0] + "..."
	send(&page_lst)
	http.Redirect(w, r, "/index", http.StatusFound)
}


//viewHandler is to view the blog
func viewHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	receive(&page_lst)
	title := r.URL.Path[len("/view/"):]
	p, _ := find_page(title, page_lst)
	path,_ := os.Getwd()
	t,_ := template.ParseFiles(path+"/Frontend/view.html")
	t.Execute(w,p)
}

//deleteHandler is to delete the blog in the edit page
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var page_lst Pages
	receive(&page_lst)
	title := r.URL.Path[len("/delete/"):]
	_,index := find_page(title,page_lst)
	page_lst.Pages = append(page_lst.Pages[:index],page_lst.Pages[index+1:]...)
	send(&page_lst)
	http.Redirect(w, r, "/index", http.StatusFound)
}

//frontPageHandler is to create the index page
func frontPageHandler(w http.ResponseWriter, r *http.Request){
	var page_lst Pages
	receive(&page_lst)
	path,_ := os.Getwd()
	t,_ := template.ParseFiles(path+"/Frontend/index.html")
	t.Execute(w,page_lst)
}

func receive(page_lst *Pages){
	fmt.Fprintf(server, "readR")
	err := decoder.Decode(page_lst)
	if err != nil{
		fmt.Println(err.Error())
	}
}

func send(page_lst *Pages) {
	fmt.Fprintf(server, "writeR")
	time.Sleep(time.Millisecond)
	err := encoder.Encode(page_lst)
	if err != nil{
		fmt.Println(err.Error())
	}
}

func testAlive(){
	for {
		time.Sleep(time.Second)
		server.SetReadDeadline(time.Now().Add(time.Second*3))
		fmt.Fprintf(server, "HiR")
		_, err := bufio.NewReader(server).ReadString('R')
		if err != nil{
			loc, _ := time.LoadLocation("UTC")
			now := time.Now().In(loc)
			fmt.Println("Detected failure on localhost:" + *portPtr+ " at " + now.String())
		}
	}
}

func main() {
	portPtr = flag.String("listen", "8080", "a string")
	serverPtr := flag.String("backend", "8090", "a string")
	flag.Parse()
	server, _= net.Dial("tcp", "localhost:"+*serverPtr)
	encoder = gob.NewEncoder(server)
	decoder = gob.NewDecoder(server)
	go testAlive()
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/index", frontPageHandler)
	fmt.Println("Go to 127.0.0.1:"+*portPtr+"/index")
	log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
}