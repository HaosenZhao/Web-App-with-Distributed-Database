Haosen Zhao (hz1126)

• Instructions: Set the GOPATH and cd to this directory (same as instructions in phase 3)

Backend: 
1. cd to Backend directory (important! or path error may happen) and build back.go
For example: go build back.go, which will create an executable file called back
2. Run the file, two command arguments provided:
a. —-listen, which will indelicate the address of this server, by default is 8090 in local
b. —-backend, which will indicate the address of other backends
for example: $ ./back --listen 8090 --backend :8091,:8092 will start the current backend on 8090

Frontend: 
1.cd to Frontend directory (important! or path error may happen)and build front.go
For example: go build front.go, which will create an executable file called front

2.Run the file, 2 command-line arguments provided
a. --listen, which will indicate the port number to accept HTTP connections on. If unspecified, application listen on port 8080.
For example, ./front —listen 8999 will cause it to listen on TCP port 8999
a link will show if user runs the file, user can just copy and paste that to a browser.

b. --backend, which indicate the endpoint used to communicate with the backend. 
An endpoint consists of the the addresses of all backends separated by commas.
if the hostname is omitted, it is assumed to refer to the local host. 
for example:
./front --listen 8080 --backend :8090,:8091,:8092 will start the frontend listening to port 8080, and using local:8090,8091,8092 as backends

3. On website, Your Blog shows the blog created with the title and the first sentence in the body
(the default is A Gallry Called Hope). If user click on the title, this will lead to the view page.
In view page, user can choose to edit, delete, or go back to main page.
In edit page, user can edit or update the blog.
If user click on delete, this will delete current page and return to main page.
Create a new blog here allows the user to create an empty blog with the title given, and lead to the edit page.
If the number of working replicas falls below quorum, your system may reject requests until the quorum is recovered.

• The state of your work: 
1. The backend now can handle concurrent requests with distributed backends
2. The system can withstand the failure of up to lower(n/2) failure of backends 
3. -If some number of backends start, the data will be entered through the web interface. 
   -Then one or more of the back ends are forcibly terminated, the web interface still remain responsive if quorum remains.
   -If the number of working replicas falls below quorum, the system will reject requests until the quorum is recovered, all operations will lead to the error page
   -If some backends reconnect and form a quorum, the system will be responsive again and data previously entered are available
4. The backend will have a id which is a random number between 1-1024, and if the backend is lost in connection, the front end will pop up error message saying lost in connection to this id, but the functionalities will not be influenced as long as it has a quorum of backend
5. Sometimes there are network partition issue but will be recovered in the heartbeat


•  Any other resources used: 
The sample code in official website of The Go Programming Language, In Writing Web Application: https://golang.org/doc/articles/wiki/#tmp_13
Golang example on handling TCP connection and setting timeout. In GitHubGist: https://gist.github.com/hongster/04660a20f2498fb7b680


•  Important design decisions made in completing this assignment: 
I used a distributed hash map to create the distributed system. Each be distributed with a id between 1-1024 and the title of each blog would be hashed and store in the backend with closet id to the hash number.
I use this strategy since the contents of the system are blogs with only title and content, it’s easy for the frontend to read from the backends and gather all blogs. And the communication between backends would be minimized, no need to worry about the leader of the backends.
However, it also have disadvantage concerning about the add and remove backends, since there will be data transfer between the backends. Also, if the view could be changed, which is not a problem in this system since the backend should be specified initially, adding new backend into the view would be a problem.
Actually, I think the system could withstand the failure to just one backend existed, but the quorum should still be used for the requirements.


• Any additional thoughts or questions about the assignment: 
This assignment is interesting and instructions are clear. The balance between consistency and availability was always in my mind when I do the phase 4, and I need to consider about different cases that would undermine the system, which was not an easy thing to do.
