Haosen Zhao (hz1126)
• Instructions: 1. Set the GOPATH and cd to this directory, then build the blog.go.
                For example: go build blog.go
                this will create an executable file called blog
                2. Run the file, the default port of the web app is 8080.
                If any port number after listen command is given, the port will be that number
                For example: ./blog will create the web app on localhost:8080
                             ./blog --listen 1234 will create the web app on localhost:1234
                a link will show if user runs the file, user can just copy and paste that to a browser.
                3. On website, Your Blog shows the blog created with the title and the first sentence in the body
                (the default is A Gallry Called Hope). If user click on the title, this will lead to the view page.
                In view page, user can choose to edit, delete, or go back to main page.
                In edit page, user can edit or update the blog.
                If user click on delete, this will delete current page and return to main page.
                Create a new blog here allows the user to create an empty blog with the title given, and lead to
                the edit page.

• The state of your work: complete all the functions required, Create, Read, Update, and Delete
                          optional command-line argument --listen works
                          application’s data in stored in memory, in a global variable
                          no database to store, no file in disk except for default blog
                          each time start application, the data set is identical with default blog

•  Any other resources used: The sample code in official website of The Go Programming Language,
                             In Writing Web Application: https://golang.org/doc/articles/wiki/#tmp_13

• Any additional thoughts or questions about the assignment: this assignment is interesting, and I learn how to use
  template, http and flag in Golang and how to create a simple web app.