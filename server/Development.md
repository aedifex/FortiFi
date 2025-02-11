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
├── pkg/ -- reusable structs
│   └── utils/
│       └── jwt.go
│
├── configs/
│   └── config.yaml
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

Connecting to Firebase allows for notifications to be sent to registered devices. There are a few steps to setting this firebase connection up in a development environment. 

Prerequisites:
- Firebase account [here](https://firebase.google.com/)
- XCode
- Ios Device
- *The notifications might not be able to run without an Apple Developer Account*

Create a new Firebase IOS project and navigate to the project console. Click the gear icon, then project settings. Under the cloud messaging tab, we can add the Apple Push Notifications Key to interact with the push notifications service. In APNs authentication key under iOS app configuration, click the Upload button. 

If you are following this setup and wish to test notifications locally, contact [@Jonathan](jonathan.nguyen@berkeley.edu) for APN keys. Browse to the location where you saved your key, select it, and click Open. Add the key ID for the key and click Upload. 

Start the Xcode project with a <b>physical</b> IOS device plugged in and selected as the target build machine. After this, follow the API guidelines on creating a user, initializing the PI, and sending notifications.

### Running the server

To run the server, use `make local-dev`. The server can be queried via http using curl, postman, or other methods. To reset the database, run `make clean-database`. To empty the logs, run `make clean-logs`. To clean both, run `make clean` (this will only work if both the database and the logs exist). To recreate the database after wiping, run `make database-dev`. Enter your `mysql` password when prompted.