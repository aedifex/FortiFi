# FortiFi Auth API

This api serves as authentication services for users and hardware devices wishing to interact with the FortiFi product. The follow document specifies the routes and flow of devices interacting with this api along with the appropriate request headers and bodies.

## Paths

### Routes for Pi Devices
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
            "id": id
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
response_headers:
    - Jwt: short-lived token
    - Refresh: long-lived token
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
        - id: string
        - details: string
        - ts: timestamp string
        - expires: timestamp string
        - type: string (1 for port scan, 2 for ddos)
        - src: string
        - dst: string
        - confidence: int
    - example:
        {
            "event": {
                "id": "id",
                "details": "there has been an intrusion on your network",
                "ts": "2006-01-02 15:04:05",
                "expires": "2006-01-02 15:04:05",
                "type": "1",
                "src": "10.0.1.1",
                "dst": "10.0.1.2",
                "confidence": 100
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
    - 200: OK

response_body: []
response_headers: []
```

<!-- Update Weekly Distribution -->
<b>Update Weekly Distribution</b>

```yaml
path: /UpdateWeeklyDistribution
description: Update the weekly distribution of eventsfor a user

methods:
    - POST

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - benign: int
    - port_scan: int
    - ddos: int
    - example:
        {
            "benign": 124,
            "port_scan": 13,
            "ddos": 2
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 400: bad request
        fix: check request body follows format above and all fields are present
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers: []
```

<!-- Reset Weekly Distribution -->
<b>Reset Weekly Distribution</b>

```yaml
path: /ResetWeeklyDistribution
description: Reset the weekly distribution of events for a user

methods:
    - POST

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - week_total: int
    - example:
        {
            "week_total": 124
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers: []
```

<!-- Add Device -->
<b>Add Device</b>

```yaml
path: /AddDevice
description: Add a device to a user

methods:
    - POST

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - name: string
    - ip_address: string
    - mac_address: string
    - example:
        {
            "name": "smartTV",
            "ip_address": "10.0.1.1",
            "mac_address": "00:00:00:00:00:00"
        }

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 400: bad request
        fix: check request body follows format above and all fields are present
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: []
response_headers: []
```

### Routes for Client Devices
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
                "id": "id",
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
    - 404: the user is not found
        fix: check the provided email belongs to an account
    - 401: unauthorized
        fix: check the password is correct
    - 500: internal server error
        fix: check server logs
    - 201: OK

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

<!-- Get Events  -->
<b>Get Events</b>

```yaml
path: /GetUserEvents
description: Get anamolous or threat events for a user

methods:
    - GET

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - events: [Events]
        - id: string
        - details: string
        - ts: timestamp string
        - expires: timestamp string
        - type: string (1 for port scan, 2 for ddos)
        - src_ip: string
        - dst_ip: string
        
response_headers: []
```

<!-- Get Weekly Distribution -->
<b>Get Weekly Distribution</b>

```yaml
path: /GetWeeklyDistribution
description: Get the weekly distribution of events for a user

methods:
    - GET

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - normal: int
    - anomalous: int
    - malicious: int
    - example:
        {
            "normal": 124,
            "anomalous": 13,
            "malicious": 2
        }

response_headers: []
```

<!-- Get Devices -->
<b>Get Devices</b>

```yaml
path: /GetDevices
description: Get all devices for a user

methods:
    - GET

query_params: []

headers:
    - Authorization: Bearer <jwt token>

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - devices: [Devices]
        - id: int
        - name: string
        - ip_address: string
        - mac_address: string
    - example:
        {
            "devices": [
                {
                    "id": 1,
                    "name": "smartTV",
                    "ip_address": "10.0.1.1",
                    "mac_address": "00:00:00:00:00:00"
                }
            ]
        }

response_headers: []
```

<!-- Get Threat Assistance -->
<b>Get Threat Assistance</b>

```yaml
path: /GetThreatAssistance
description: Get threat assistance for a user

methods:
    - GET

query_params:
    - threatId

headers:
    - Authorization: Bearer <jwt token>

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database   
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - response: string

```

<!-- Get More Assistance -->
<b>Get More Assistance</b>

```yaml
path: /GetMoreAssistance
description: Get more assistance for a user

methods:
    - GET

query_params:
    - threatId

headers:
    - Authorization: Bearer <jwt token>

request_body: json
    - query: string

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - response: string

response_headers: []
```

<!-- Get Recommendations -->
<b>Get Recommendations</b>

```yaml
path: /GetRecommendations
description: Get recommendations for a user

methods:
    - GET

query_params:
    - threatId

headers:
    - Authorization: Bearer <jwt token>

request_body: []

responses:
    - 405: method not allowed
        fix: check http method
    - 401: unauthorized
        fix: check the jwt header and ensure valid
    - 404: not found
        fix: check the user entry in database
    - 500: internal server error
        fix: check server logs
    - 200: OK

response_body: json
    - response: string

response_headers: []
```


