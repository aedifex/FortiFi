# Development Contribution Guide

This file contains information regarding contributing to the server code

--- 

## File Structure
```sh
/server
│
├── cmd/
│   └── server/
│       └── main.go -- server initialization and run
│
├── internal/
│   ├── handler/
│   │   └── user_handler.go -- routes handling
│   ├── database/
│   │   └── database.go -- database calls, structs, logic
│   └── middleware/
│       └── auth.go -- auth middleware
│
├── pkg/ -- reusable structs (loggers etc)
│   └── utils/
│       └── logger.go
│
├── api/
│   └── user_api.proto
│
├── configs/
│   └── config.yaml
│
├── web/
│   └── index.html
│
├── scripts/
│   └── setup.sh
│
├── go.mod
├── go.sum
├── README.md
```

## Development Running and Testing

### Database Configuration
To run the server locally, make sure to have the following steps completed:

1. Download and install `mysql` for your system. 
    - Linux
        ```sh
        sudo apt update
        sudo apt install mysql-server
        ```
    - Mac
        ```sh
        brew install mysql
        ```
    - [Windows Download](https://dev.mysql.com/downloads/installer/)

2. Ensure `mysql` is running
    - Linux
        ```sh
        sudo systemctl start mysql
        ```
    - Mac
        ```sh
        brew services start mysql
        ```
    - Windows
        - Open the Services app (services.msc).
        - Find the "MySQL" service, right-click it, and click "Start."
        
3. When the server is run (later steps), the database will be properly setup for endpoint testing.

### Environment configurations

For running in development environment, add the following values to `config/dev.config.yaml`
```yaml
port: 
  ":3000"
db_user:
  "<mysql_username(most likely root)>"
db_pass:
  "<mysql_password>"
db_url:
  "127.0.0.1:3306"
db_name:
  "FortiFi"
signing_key:
  "b2e138d8553ea7d7ff8731e87e41406277bd4c98"
```

### Running the server

To run the server, use `make local-dev`. The server can be queried via http using curl, postman, or other methods. To reset the database, run `make clean`. Enter your `mysql` password when prompted.

## Endpoints

*[POST] /NotifyIntrusion*

```json
to be implemented
```
*[POST] /CreateUser*

```json
{
    "Id": string, // This should be a unique id from the raspberry pi
    "first_name": string,
    "last_name": string,
    "email": string,
    "password": string
}
```

returns `201 CREATED` on success. 

*/Login*

```json
{
    "email": string,
    "password": string
}
```

returns `200 OK` on with valid JWT and refresh token in the header on success.