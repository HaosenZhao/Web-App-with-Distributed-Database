package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

//page_lst is a global variable of Pages
var page_lst = Pages{Pages: []*Page{}}
//templates is the template for edit and view
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

//initialize function initialize the web main page and insert the default blog, which is saved in a text file
//and all other data saved in local
func initialize(){
	initial_page := Page{}
	initial_page.Title = "A Gallry Called Hope"
	body,_ := ioutil.ReadFile("A Gallry Called Hope.txt")
	body_str := string(body)
	initial_page.Body = body_str
	first_string := strings.SplitAfter(body_str, ".")[0] + "..."
	initial_page.FirstLine = first_string
	page_lst.Pages = append(page_lst.Pages, &initial_page)
}

//find_page is to return the pointer and index of page using a title
func find_page(title string)(*Page, int){
	for i := 0; i < len(page_lst.Pages); i++{
		if page_lst.Pages[i].Title == title{
			return page_lst.Pages[i],i
		}
	}
	return nil,999
}

//renderTemplate is to execute the template created
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html",p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//createHandler is to create a new blog in the main apge
func createHandler(w http.ResponseWriter, r *http.Request){
	title := r.FormValue("title")
	pag,_ := find_page(title)
	counter:=0
	ori_title := title
	for pag != nil{
		counter++
		title = ori_title+"("+strconv.Itoa(counter)+")"
		pag,_ = find_page(title)
	}
	http.Redirect(w, r, "/edit/"+title, http.StatusFound)
}

//editHandler is to edit the page created
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p,_ := find_page(title)
	if p == nil{
		p =  &Page{Title:title, Body:""}
		page_lst.Pages = append(page_lst.Pages, p)
	}
	renderTemplate(w, "edit", p)
}

//saveHandler is to save the change in blog after edit
func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p, _ := find_page(title)
	p.Body = body
	p.FirstLine = strings.SplitAfter(body, ".")[0] + "..."
	http.Redirect(w, r, "/", http.StatusFound)
}


//viewHandler is to view the blog
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := find_page(title)
	renderTemplate(w, "view", p)
}

//deleteHandler is to delete the blog in the edit page
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/delete/"):]
	_,index := find_page(title)
	page_lst.Pages = append(page_lst.Pages[:index],page_lst.Pages[index+1:]...)
	http.Redirect(w, r, "/", http.StatusFound)
}

//frontPageHandler is to create the index page
func frontPageHandler(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("index.html")
	t.Execute(w,page_lst)
}

func main() {
	portPtr := flag.String("listen", "8080", "a string")
	flag.Parse()
	initialize()
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/", frontPageHandler)
	http.HandleFunc("/create", createHandler)
	fmt.Println("Go to 127.0.0.1:"+*portPtr)
	log.Fatal(http.ListenAndServe(":"+*portPtr, nil))
}