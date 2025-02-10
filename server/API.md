# FortiFi Auth API

This api serves as authentication services for users and hardware devices wishing to interact with the FortiFi product. The follow document specifies the routes and flow of devices interacting with this api along with the appropriate request headers and bodies.

## Paths

<!-- PiInit Path -->
<b>Pi Initialization</b>

```yaml
path: /PiInit
description: Get auth tokens for future requests

methods:
    - POST

query_params: []

headers: []

request_body: json
    - id: string (should be a uuid associated with the pi)
    - example:
        {
            "id": "userId123"
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 400: bad request
        fix: check request body -- ensure id is there
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers:
    - Jwt: short-lived token
    - Refresh: long-lived token
```

<!-- Pi Refresh Path -->
<b>Pi Refresh</b>

```yaml
path: /RefreshPi
description: Get auth tokens for future requests

methods:
    - GET

query_params:
    - id

headers:
    - Refresh: refresh token from previous refresh/init request

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the refresh token
    - 404: not found
        fix: check the id is correct
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers: []
```

<!-- User Create Path -->
<b>Create User</b>

```yaml
path: /CreateUser
description: register a new user

methods:
    - POST

query_params: []

headers: []

request_body: json
    - user: user information
        - id: string
        - first_name: string
        - last_name: string
        - email: string
        - password: string
    - example:
        {
            "user": {
                "id": "userId123",
                "first_name":"oski",
                "last_name":"bear",
                "email":"oski@berkeley.edu",
                "password":"go bears 2025"
            }
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 400: bad request
        fix: check request body follows format above and all fields are present
    - 409: the user already exists
    - 500: internal server error
        fix: check server logs
    - 201: CREATED

response_body: []
response_headers: []
```

<!-- User Login Path -->
<b>Login User</b>

```yaml
path: /Login
description: login a new user

methods:
    - POST

query_params: []

headers: []

request_body: json
    - user: user information
        - email: string
        - password: string
    - example:
        {
            "user": {
                "email":"oski@berkeley.edu",
                "password":"go bears 2025"
            }
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 400: bad request
        fix: check request body follows format above and all fields are present
    - 409: the user already exists
    - 500: internal server error
        fix: check server logs
    - 201: CREATED

response_body: []
response_headers:
    - Jwt: short-lived jwt token
    - Refresh: long-lived refresh token
```

<!-- User Refresh Path -->
<b>User Refresh</b>

```yaml
path: /UserRefresh
description: Get auth tokens for future requests

methods:
    - GET

query_params:
    - id

headers:
    - Refresh: refresh token from previous refresh/login request

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the refresh token
    - 404: not found
        fix: check the id is correct
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers:
    - Jwt: short-lived jwt token
    - Refresh: long-lived refresh token
```

<!-- User Notification Token Path -->
<b>Set User Notifications Token</b>

```yaml
path: /UpdateFcm
description: Set device notifications token

methods:
    - POST

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - fcm_token: string (can get from xcode console)

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 400: bad request
        fix: check fcm_token is in body as json
    - 404: not found
        fix: check the jwt is correct -- id is not found in db
    - 500: internal server error
        fix: check server logs
    - 202: Accepted

response_body: []
response_headers:
    - Jwt: short-lived jwt token
    - Refresh: long-lived refresh token
```

<!-- Notifications Path -->
<b>Send Notifications</b>

```yaml
path: /NotifyIntrusion
description: Send a notification to user

methods:
    - POST

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - event:
        - details: string
        - ts: timestamp string
        - expires: timestamp string
    - example:
        {
            "event": {
                "details": "there has been an intrusion on your network",
                "ts": "2006-01-02 15:04:05",
                "expires": "2006-01-02 15:04:05"
            }
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 400: bad request
        fix: check event is in body as json and all event fields are present
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 202: Accepted

response_body: []
response_headers:
    - Jwt: short-lived jwt token
    - Refresh: long-lived refresh token
```