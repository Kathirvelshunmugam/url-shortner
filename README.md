# url-shortner
This is CRUD operations on MongoDB written in Golang. You can create, read, update, and delete URLs from the MongoDB instance using HTTP requests using Postman API.

**How to Run?**

To start our project we will install our dependencies in our case will be the MongoDB go driver, gin package and shortid package:
```
go get go.mongodb.org/mongo-driver
go get github.com/gin-gonic/gin
go get github.com/teris-io/shortid
```
Following the above next install the dependencies using:
go get ./...

Finally, run the app on port 8080:
go run .

**Endpoints:**
```
GET    /:code   
POST   /shorten 
PUT    /:id    
DELETE /:id  

**Get URL**
```
This endpoint retrieves an entire site of a given URL.
Send a GET request to /:code:

Below is an example, Using Postman API GET send this request along with code generated in MongoDB
```
GET - http://localhost:8080/IJbrJ024R
```
Response:

You will get the entire site in multiple views in the response console of Postman API

**Create URL**

This endpoint inserts a URL in the url collection of the Project database.

Send a POST request to /shorten:

Below is an example, Using Postman API POST this request along with adding an URL in the Body by changing the file type(raw) into JSON
```
POST - http://localhost:8080/shorten
```
Body:
```
{
    "long_url":"https://www.youtube.com/watch?v=-LEye3S"
}
```
Response:
```
{
    "db_id": "63be8d7f13e5f9f6dba1dda4",
    "newUrl": "http://localhost:8080/fkywCbh4R"
}
```
**Update URL:**

The below endpoint helps to update the existing URL in the database.

Send PUT request to /:id
```
PUT - http://localhost:8080/f5fd44e3-6434-4949-995e-e2eccf69c578
```

Response:

{
    "error": false,
    "message": "Long Url Updated"
}

**Delete URL:**

The below endpoint performs a delete URL operation in the database

Send DELETE request to /:id
```
DELETE - http://localhost:8080/f5fd44e3-6434-4949-995e-e2eccf69c578
```
Response:

{
    "error": false,
    "message": "Long Url Deleted"
}

The above procedure helps you to perform URL CRUD operations in MongoDB.

Note:

:code and :id that has been used above depend on my MongoDB creation. That is unique for every Developer.
