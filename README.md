# RestMail #

### What is this repository for? ###

* REST API written in Go. Available three endpoints, described in lower part of README. The application has connection with Cassandra, where messages are stored. App allows to send messages as emails.
* 1.0

### How do I get set up? ###

* To send emails, you have to clone this repository and change credentials to your own smtp, in code. 
* Database is created by following commands
>CREATE KEYSPACE IF NOT EXISTS messages WITH replication = {'class' : 'SimpleStrategy', 'replication_factor' : 1};

>CREATE TABLE messages.messages (email text, title text, content text, magic_number int, PRIMARY KEY (email, magic_number));
* To run this app, you have to run it in docker container on 8080 port. It also requires cassandra in docker container on port 9042. 
* Command to run cassandra docker docker run -p 9042:9042 -d --name cassandra cassandra

### API Endpoints ###

* POST localhost:8080/api/message - add email to storage
* POST localhost:8080/api/send - send emails with given magic_number and remove them from storage
* GET localhost:8080/api/messages/{emailValue} - return all messages with given email

### API usage examples ###

* CURL -X POST localhost:8080/api/message -d '{"email":"test@test.com","title":"Email title","content":"There is message","magic_number":10}'
* CURL -X POST localhost:8080/api/send -d '{"magic_number":10}'
* CURL -X GET localhost:8080/api/messages/test@test.com
