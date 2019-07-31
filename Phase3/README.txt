Haosen Zhao (hz1126)

• Instructions: Set the GOPATH and cd to this directory (same as instructions in phase 2)

Backend: 
1. cd to Backend directory (important! or path error may happen) and build back.go
For example: go build back.go, which will create an executable file called back
2. Run the file, the default port of the backend is 8090.
If any port number after listen command is given, the port will be that number
For example: ./blog will create the web app on localhost:8090
./blog --listen 1234 will create the backend listening to port 1234

Frontend: 
1.cd to Frontend directory (important! or path error may happen)and build front.go
For example: go build front.go, which will create an executable file called front

2.Run the file, 2 command-line arguments provided
a. --listen, which will indicate the port number to accept HTTP connections on. 
For example, ./front —listen 8999 will cause it to listen on TCP port 8999
If unspecified, application listen on port 8080.
a link will show if user runs the file, user can just copy and paste that to a browser.

b. --backend, which indicate the endpoint used to communicate with the backend. 
An endpoint consists of a hostname and port, separated by a colon; 
if the hostname is omitted, it is assumed to refer to the local host. 
For example,  --backend something.com:9000 will cause it connect to host something.com on port 9000. 
As another example, —backend :9100 will cause it to connect to localhost at port 9100. 
If unspecified, application connect to the back end at localhost:8090.

3. On website, Your Blog shows the blog created with the title and the first sentence in the body
(the default is A Gallry Called Hope). If user click on the title, this will lead to the view page.
In view page, user can choose to edit, delete, or go back to main page.
In edit page, user can edit or update the blog.
If user click on delete, this will delete current page and return to main page.
Create a new blog here allows the user to create an empty blog with the title given, and lead to the edit page.


• The state of your work: 
1. The backend now can handle concurrent requests from front ends. No race condition, deadlock, live lock detected.
2. Test the front end website under load by using Vegeta attack, Success rate is 100% in a variety of load test.
3. Front end now have a failure detector, which will detect any disconnection or failure from the backend, the test rate is 1sec/test and waiting rate is 3 sec/test. 
4. Synchronization and locking approach is constructed in this phase. 
5. The failure detector is isolated from other web handlers and backend connection, which is implemented in a thread.


•  Any other resources used: 
The sample code in official website of The Go Programming Language, In Writing Web Application: https://golang.org/doc/articles/wiki/#tmp_13
Golang example on handling TCP connection and setting timeout. In GitHubGist: https://gist.github.com/hongster/04660a20f2498fb7b680


•  Important design decisions made in completing this assignment: 
Currently I wrapped up the entire database in backend for synchronization. I also use a mutex for every read/write to the database for safety. There is no safety issue detected, since this mechanism is simple, and safe to use, but has a large granularity. And this mechanism also provides great consistency to all the front ends, since the front end will get all the data from the backend when they makes a required call. I use this approach since the web app is a blog with only text data, there is nothing large to transmit, and it’s ok to send all the database to the user. But there is a possible situation when two clients save in the same time, since the later always win that the former user’s data might be covered by the database saved by the later user. This is a functional disadvantage in this approach. 

I use Vetaga attack to test the under load condition in a variety of access patterns for handler. And I solved problems detected by making the request redirect to the index page if there is a possible conflict operation. And I finally got 100.00% success ratio to all the Vegeta attack tested and small latencies, the reason of which should be the small size of the data to be transmitted. I also draw a graph of a report and noticed that there is a short time when the latency became so large. I guess the reason for the fluctuation in latency might be when several delete operations of the same data happened in the same time and one made the deletion but others are required to redirect to the index page.
One copy of report and graph is included in the folder.

For failure detect approach, I use a pic-act method to detect any possible failure. I make the client to communicate with the server every second, and if there is no reply from the server in next three seconds, a error message would be printed. I could use heartbeat for failure detection but I want to try pin-act for I can set the wait time to see the balance between completion and accuracy. I also notice that if the server is died, the error message would be printed every one second instead of 4 seconds (1 sec/test + 3 waitSec/test), since the the front end would not wait for a dead server but acknowledge that there won’t be any reconnection for the server. 

• Any additional thoughts or questions about the assignment: 
This assignment is interesting and instructions are clear. Learn how to protect the data and how to use failure detect mechanism taught in the class. 
