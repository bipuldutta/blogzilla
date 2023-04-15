# Blogzilla

This is a simple blogging application built using Go and Postgres. This application is an attempt to follow
The Clean Architecture as I have used it in the past and like the simplicity & testability of each layer in the 
project. More about it can be found https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html 

This service provides REST APIs for creating, updating, reading, and deleting blog posts.
And also have feature to register users, create and manage roles with various permissions.

## Special Note

The original version of this project was almost built in like ~5 hours for me to refresh 
my knowledge of building a backend product from scratch. Please use this as your 
learning reference and nothing more. Hope it helps.

## Requirements

- Go (any latest version like 1.19.x or above)
- Docker (Engine 20.10.x). Use it to host a Postgres container (You don't really need it 
  if you have an accessible Postgres)
- Docker Compose (v2.15). If you are using Docker for the postgres. My plan is to use it in future to add 
  more containers (maybe prometheus to add metrics, or redis to cache some important data that I do not know yet)  
- macOS (other OS can be used but the documentation could be different)

## How to Build

At first clone the repository as follows:

```
git clone git@github.com/bipuldutta/blogzilla
```

Go to the directory where you have cloned the project
```
>cd server
>go build
```
This will build a macOS executable name `server`.
Note: you may see package dependency errors, if you do then run `go get <package>` command to get them.  

## How to Run

Make sure that a postgres database is running with the `postgres` configuration found in 
the `config/config-local.yml` file. When the server is started it will take care of creating 
the database tables, roles, and the admin user defined in the `defaultuser` section
in the `config/config-local.yml` file. To run the server buil

```
>./server
```

### Roles
Following roles are available
- **admin**: the administrators of the system.
- **editor**: user who can create, update, delete blogs.
- **viewer**: blog viewer.

### Authentication and Authorization

We are using JWT (https://jwt.io/introduction) as the result of a successful user authentication and use it for subsequent
calls to protected endpoints. JWT is a flexible token format where we could embed other important information 
like permissions, user id, expiration, etc. for the protected endpoints to be able to validate without requiring additional
interactions with some other Auth service. Main goal is to keep it simple and still be robust.

## How to Test

Following are several APIs can be tested

### Registration

This endpoint will let you create a new user. By default, a new user gets the `editor` and `viewer` roles.
In the future, we can have more endpoints to assign/unassign roles. Or add new roles with different accessibility.

Request:
```
curl --request POST \
  --url http://localhost:8080/v1/register \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "example_user123",
	"password": "Pa$$w0rd2023!",
	"firstName": "James",
	"lastName": "Parker"
}'
```
Response:`201 Created`

### Login

User can login using the above username and password
Request:
```
curl --request POST \
  --url http://localhost:8080/v1/login \
  --header 'Content-Type: application/json' \
  --data '{
	"username": "example_user123",
	"password": "Pa$$w0rd2023!"
}'
```
Response:
```
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJwZXJtaXNzaW9ucyI6bnVsbCwiZXhwIjoxNjgxMzY1ODQzLCJpYXQiOjE2ODEzNjQ2NDN9.VriAJR8whuXkkE4M2FmOyoaPFbl-PpWjJwJuEbqFejo"}
```
The `token` is a JWT token which includes user id, permissions, and expiration time claims. 
After login, all other endpoints will require this as a `Bearer` token in the `Athorization` header. 
Following examples will show how it is passed to the endpoints.

### Create a blog

Request:
```
curl --request POST \
  --url http://localhost:8080/v1/blogs \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJwZXJtaXNzaW9ucyI6WyJkZWxldGVfdXNlciIsImNyZWF0ZV9ibG9nIiwidXBkYXRlX2Jsb2ciLCJkZWxldGVfYmxvZyIsImNyZWF0ZV91c2VyIiwicmVhZF91c2VyIiwidXBkYXRlX3VzZXIiLCJyZWFkX2Jsb2ciXSwiZXhwIjoxNjgxMzY1MTcyLCJpYXQiOjE2ODEzNjM5NzJ9.VNS_wN6yjXVzhQe9jg3Ml4PZ2LeVrWYaESSy7Nhia2I' \
  --header 'Content-Type: application/json' \
  --data '{
    "title": "My First Blog Post",
    "content": "Lorem ipsum dolor sit amet",
    "tags": "foo,bar"
}'
```
Response:
```
{"id":2}
```

### Search Blogs

Request:
```
curl --request GET \
  --url 'http://localhost:8080/v1/blogs?q=foo&offset=0&limit=10' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjMsIlBlcm1pc3Npb25zIjpbImNyZWF0ZV9ibG9nIiwidXBkYXRlX2Jsb2ciLCJkZWxldGVfYmxvZyIsInJlYWRfdXNlciIsInJlYWRfYmxvZyJdLCJleHAiOjE2ODE1MjM2MDksImlhdCI6MTY4MTUyMjQwOX0._5hn_QHZbr3r-3FfXWgJNDyZCEUWZuIC1vKErYpLggg' \
  --header 'Content-Type: application/json'
```
Response:
```
[
  {
    "ID": 1,
    "UserID": 3,
    "Title": "My Second Blog Post",
    "Content": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed vel erat ultricies, vulputate leo a, malesuada eros. Sed euismod tortor vitae nisl blandit, quis bibendum sapien ullamcorper. Proin luctus mauris eu enim finibus, non convallis risus consectetur. Sed pulvinar, nunc non consectetur bibendum, velit arcu vestibulum massa, vitae faucibus velit magna ac turpis.",
    "Tags": "mindfulnes,foo",
    "CreatedAt": "2023-04-15T01:33:56.37797Z",
    "UpdatedAt": "2023-04-15T01:33:56.37797Z"
  }
]
```