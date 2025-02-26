package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aedifex/FortiFi/config"
	"github.com/aedifex/FortiFi/internal/database"
	"github.com/aedifex/FortiFi/internal/requests"
)

var (
	ReqFail = formatError("failed to build request")

	BadStatus = func(exp int, got int) string {
		return formatError("bad status")(fmt.Errorf("expected %d but got %d ", exp, got))
	}

	BadResponseHeader = "invalid response: expected headers"

	InvalidBody = func(exp interface{}, got interface{}) string {
		return formatError("invalid body")(fmt.Errorf("\nexpected: %v\ngot %v", exp, got))
	}
)

var (
	piJwt       = ""
	piRefresh   = ""
	userJwt     = ""
	userRefresh = ""

	// these are used to test invalidated tokens after a refresh request goes through
	jwtTemp     = ""
	refreshTemp = ""

	id        = "userId123"
	firstName = "Oski"
	lastName  = "Bear"
	email     = "oskibear@berkeley.edu"
	password  = "Go Bears!"
	fcmToken  = "c9g9Wdmp90NchwYMhiiLBp:APA91bHfeKLGu921KeSj45uikQbhg1_Gx44qBjHErrjDoMzSIag5fGJdUotOCOjCQumLv2etUbe_e_gfJNKOIQEhUua6KIcp7zQcgGetkleiWPLJgRJ3GcY"
)

var server = setupTestServer()
var testHandler = server.httpServer.Handler

type testCase struct {
	name          string
	correctStatus int
	requestBody   interface{}
	getsTokens    bool
	jwt           string
	refresh       string
	setup         func(req *http.Request)
	jwtTarget     *string
	refreshTarget *string
	expectedBody  interface{}
}

func formatError(s string) func(err error) string {
	return func(err error) string {
		return fmt.Sprintf("%s: %s", s, err.Error())
	}
}

func setupTestServer() *fortifiServer {
	testing.Init()

	// Setup environment
	config := &config.Config{
		Port:        ":3000",
		DB_USER:     "root",
		DB_PASS:     "",
		DB_URL:      "localhost:3306",
		DB_NAME:     "FortiFi",
		SIGNING_KEY: "b2e138d8553ea7d7ff8731e87e41406277bd4c98",
		FcmKeyPath:  "/harness/firebase_auth.json",
	}

	// Create new FortifiServer
	server := newServer(config)
	if server == nil {
		log.Fatalf("Server was not initialized correctly")
	}

	return server
}

func Marshal(body interface{}, t *testing.T) io.Reader {
	t.Helper()
	if body == nil {
		return nil
	}
	marshaledBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(marshaledBody)
}

func methodTest(t *testing.T, path string) {
	t.Helper()
	t.Run("method not allowed test", func(t *testing.T) {
		req := buildRequest(t, http.MethodTrace, path, nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", piJwt))
		resp := httptest.NewRecorder()
		testHandler.ServeHTTP(resp, req)
		if resp.Code != http.StatusMethodNotAllowed {
			t.Fatalf(BadStatus(http.StatusMethodNotAllowed, resp.Code))
		}
	})
}

func missingBodyTest(t *testing.T, method string, path string) {
	t.Helper()
	t.Run("missing body test", func(t *testing.T) {
		req := buildRequest(t, method, path, nil)
		// This is to test a missing body on protected routes
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userJwt))
		resp := httptest.NewRecorder()
		testHandler.ServeHTTP(resp, req)
		if resp.Code != http.StatusBadRequest {
			t.Fatal(BadStatus(http.StatusBadRequest, resp.Code))
		}
	})
}

func buildRequest(t *testing.T, method string, path string, body interface{}) *http.Request {
	t.Helper()
	marshaledBody := Marshal(body, t)
	req, err := http.NewRequest(method, path, marshaledBody)
	if err != nil {
		t.Fatal(ReqFail(err))
	}
	return req
}

func buildTest(tc testCase, method string, path string) func(t *testing.T) {
	return func(t *testing.T) {
		methodTest(t, path)
		req := buildRequest(t, method, path, tc.requestBody)
		if tc.jwt != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tc.jwt))
		}
		if tc.refresh != "" {
			req.Header.Add("Refresh", tc.refresh)
		}
		if tc.setup != nil {
			tc.setup(req)
		}
		resp := httptest.NewRecorder()
		testHandler.ServeHTTP(resp, req)
		if resp.Code != tc.correctStatus {
			t.Fatal(BadStatus(tc.correctStatus, resp.Code))
		}
		if tc.getsTokens {
			if resp.Header().Get("jwt") == "" || resp.Header().Get("refresh") == "" {
				t.Fatal(BadResponseHeader)
			}
			jwtTemp = *tc.jwtTarget
			refreshTemp = *tc.refreshTarget
			*tc.jwtTarget = resp.Header().Get("jwt")
			*tc.refreshTarget = resp.Header().Get("refresh")
		}
		if tc.expectedBody != nil {
			// Convert both expected and got to JSON strings for comparison
			expectedJSON := Marshal(tc.expectedBody, t)
			expectedBytes, err := io.ReadAll(expectedJSON)
			if err != nil {
				t.Fatalf("failed to read expected body: %s", err)
			}
			expectedString := strings.TrimSpace(string(expectedBytes))

			got := strings.TrimSpace(resp.Body.String())

			if expectedString != got {
				t.Fatalf("Strings differ:\nexpected (len=%d): %q\ngot (len=%d): %q",
					len(expectedString), expectedString,
					len(got), got)
				t.Fatalf(InvalidBody(expectedString, got))
			}
		}
	}
}

func TestPiInit(t *testing.T) {
	path := "/PiInit"
	validMethod := http.MethodPost
	testCases := []testCase{
		{
			name:          "correct body",
			correctStatus: http.StatusOK,
			requestBody: &requests.PiInitRequest{
				Id: id,
			},
			getsTokens:    true,
			jwtTarget:     &piJwt,
			refreshTarget: &piRefresh,
		},
		{
			name:          "missing id",
			correctStatus: http.StatusBadRequest,
			requestBody:   &requests.PiInitRequest{},
			getsTokens:    false,
		},
	}

	missingBodyTest(t, validMethod, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, validMethod, path))
	}
}

func TestPiRefresh(t *testing.T) {

	path := "/RefreshPi"
	validMethod := http.MethodGet
	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			getsTokens:    true,
			refresh:       piRefresh,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
			jwtTarget:     &piJwt,
			refreshTarget: &piRefresh,
		},
		{
			name:          "retrying old tokens",
			correctStatus: http.StatusUnauthorized,
			getsTokens:    false,
			refresh:       refreshTemp,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
		},
		{
			name:          "missing id",
			correctStatus: http.StatusNotFound,
			getsTokens:    false,
			refresh:       piRefresh,
		},
		{
			name:          "invalid id",
			correctStatus: http.StatusNotFound,
			getsTokens:    false,
			refresh:       refreshTemp,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", "badId123")
				req.URL.RawQuery = query.Encode()
			},
		},
		{
			name:          "missing token",
			correctStatus: http.StatusUnauthorized,
			getsTokens:    false,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
		},
		{
			name:          "invalid token",
			correctStatus: http.StatusUnauthorized,
			getsTokens:    false,
			refresh:       "bad.refresh.token",
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, validMethod, path))
	}
}

func TestCreateUser(t *testing.T) {

	emptyBody := &requests.CreateUserRequest{}

	userCorrect := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		},
	}

	userDuplicateId := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id, // Assuming this ID is already in the database
			FirstName: firstName,
			LastName:  lastName,
			Email:     "duplicateId@berkeley.edu",
			Password:  "duplicatePassword!",
		},
	}

	userDuplicateEmail := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email, // Assuming this email is already in the database
			Password:  password,
		},
	}

	userMissingId := &requests.CreateUserRequest{
		User: &database.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		},
	}

	userMissingFirstName := &requests.CreateUserRequest{
		User: &database.User{
			Id:       id,
			LastName: lastName,
			Email:    email,
			Password: password,
		},
	}

	userMissingLastName := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id,
			FirstName: firstName,
			Email:     email,
			Password:  password,
		},
	}

	userMissingEmail := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			Password:  password,
		},
	}

	userMissingPassword := &requests.CreateUserRequest{
		User: &database.User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		},
	}

	path := "/CreateUser"
	method := http.MethodPost
	testCases := []testCase{
		{
			name:          "missing user in body",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   emptyBody,
		},
		{
			name:          "missing id field",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   userMissingId,
		},
		{
			name:          "missing first name field",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   userMissingFirstName,
		},
		{
			name:          "missing last name field",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   userMissingLastName,
		},
		{
			name:          "missing email field",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   userMissingEmail,
		},
		{
			name:          "missing password field",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   userMissingPassword,
		},
		{
			name:          "valid request",
			correctStatus: http.StatusCreated,
			getsTokens:    false,
			requestBody:   userCorrect,
		},
		{
			name:          "duplicate id field",
			correctStatus: http.StatusConflict,
			getsTokens:    false,
			requestBody:   userDuplicateId,
		},
		{
			name:          "duplicate email field",
			correctStatus: http.StatusConflict,
			getsTokens:    false,
			requestBody:   userDuplicateEmail,
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}

}

func TestLogin(t *testing.T) {

	loginMissingBody := &requests.LoginUserRequest{}

	loginCorrect := &requests.LoginUserRequest{
		User: &database.User{
			Email:    email,
			Password: password,
		},
	}

	loginMissingEmail := &requests.LoginUserRequest{
		User: &database.User{
			Password: password,
		},
	}

	loginMissingPassword := &requests.LoginUserRequest{
		User: &database.User{
			Email: email,
		},
	}

	loginBadEmail := &requests.LoginUserRequest{
		User: &database.User{
			Email:    "bad@email.com",
			Password: password,
		},
	}

	loginBadPassword := &requests.LoginUserRequest{
		User: &database.User{
			Email:    email,
			Password: "badPassword",
		},
	}

	path := "/Login"
	method := http.MethodPost
	testCases := []testCase{
		{
			name:          "missing user in body",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   loginMissingBody,
		},
		{
			name:          "missing email",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   loginMissingEmail,
		},
		{
			name:          "missing password",
			correctStatus: http.StatusBadRequest,
			getsTokens:    false,
			requestBody:   loginMissingPassword,
		},
		{
			name:          "unregistered email",
			correctStatus: http.StatusNotFound,
			getsTokens:    false,
			requestBody:   loginBadEmail,
		},
		{
			name:          "incorrect password",
			correctStatus: http.StatusUnauthorized,
			getsTokens:    false,
			requestBody:   loginBadPassword,
		},
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			getsTokens:    true,
			requestBody:   loginCorrect,
			jwtTarget:     &userJwt,
			refreshTarget: &userRefresh,
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}

}

func TestUserRefresh(t *testing.T) {

	path := "/RefreshUser"
	method := http.MethodGet
	testCases := []testCase{
		{
			name:          "valid request",
			refresh:       userRefresh,
			getsTokens:    true,
			jwtTarget:     &userJwt,
			refreshTarget: &userRefresh,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
			correctStatus: http.StatusOK,
		},
		{
			name:       "retry revoked refresh token",
			refresh:    refreshTemp,
			getsTokens: false,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "missing id",
			refresh:       userRefresh,
			getsTokens:    false,
			correctStatus: http.StatusNotFound,
		},
		{
			name:       "invalid id",
			refresh:    userRefresh,
			getsTokens: false,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", "badId123")
				req.URL.RawQuery = query.Encode()
			},
			correctStatus: http.StatusNotFound,
		},
		{
			name:       "missing token",
			getsTokens: false,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid token",
			refresh:    "bad.refresh.token",
			getsTokens: false,
			setup: func(req *http.Request) {
				query := req.URL.Query()
				query.Add("id", id)
				req.URL.RawQuery = query.Encode()
			},
			correctStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}

}

func TestUpdateFcmToken(t *testing.T) {

	method := http.MethodPost
	path := "/UpdateFcm"
	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusAccepted,
			jwt:           userJwt,
			requestBody: &requests.UpdateFcmRequest{
				FcmToken: fcmToken,
			},
		},
		{
			name:          "bad jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "bad.jwt",
			requestBody: &requests.UpdateFcmRequest{
				FcmToken: fcmToken,
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestNotifyIntrusion(t *testing.T) {
	path := "/NotifyIntrusion"
	method := http.MethodPost
	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details:    "Instrusion event details here",
					TS:         "2006-01-02 15:04:05",
					Expires:    "2026-01-02 15:04:05",
					Type:       "1",
					SrcIP:      "10.0.1.1",
					DstIP:      "10.0.1.2",
					Confidence: 100,
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing event",
			correctStatus: http.StatusBadRequest,
			requestBody:   &requests.NotifyIntrusionRequest{},
			jwt:           piJwt,
		},
		{
			name:          "missing token",
			correctStatus: http.StatusUnauthorized,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					TS:      "2006-01-02 15:04:05",
					Expires: "2026-01-02 15:04:05",
				},
			},
		},
		{
			name:          "missing details",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					TS:      "2006-01-02 15:04:05",
					Expires: "2026-01-02 15:04:05",
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing TS",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					Expires: "2026-01-02 15:04:05",
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing expires",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					TS:      "2006-01-02 15:04:05",
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing type",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					TS:      "2006-01-02 15:04:05",
					Expires: "2026-01-02 15:04:05",
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing src ip",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					TS:      "2006-01-02 15:04:05",
					Expires: "2026-01-02 15:04:05",
					Type:    "anomaly",
				},
			},
			jwt: piJwt,
		},
		{
			name:          "missing dst ip",
			correctStatus: http.StatusBadRequest,
			requestBody: &requests.NotifyIntrusionRequest{
				Event: &database.Event{
					Details: "Instrusion event details here",
					TS:      "2006-01-02 15:04:05",
					Expires: "2026-01-02 15:04:05",
					Type:    "anomaly",
				},
			},
			jwt: piJwt,
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestGetUserEvents(t *testing.T) {
	path := "/GetUserEvents"
	method := http.MethodGet

	event := map[string][]database.Event{
		"events": {
			{
				ThreatId: 1,
				Id:       id,
				Details:  "Instrusion event details here",
				TS:       "2006-01-02 15:04:05",
				Expires:  "2026-01-02 15:04:05",
				Type:     "1",
				SrcIP:    "10.0.1.1",
				DstIP:    "10.0.1.2",
			},
		},
	}

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           userJwt,
			expectedBody:  event,
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestUpdateWeeklyDistribution(t *testing.T) {
	path := "/UpdateWeeklyDistribution"
	method := http.MethodPost

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           piJwt,
			requestBody: &requests.UpdateWeeklyDistributionRequest{
				Benign:   10,
				PortScan: 5,
				DDoS:     2,
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestGetWeeklyDistribution(t *testing.T) {
	path := "/GetWeeklyDistribution"
	method := http.MethodGet

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           piJwt,
			expectedBody: &database.WeeklyDistribution{
				Benign:        10,
				PortScan:      5,
				DDoS:          2,
				PrevWeekTotal: 0,
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
		{
			name:          "user does not exist",
			correctStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestResetWeeklyDistribution(t *testing.T) {
	path := "/ResetWeeklyDistribution"
	method := http.MethodPost

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           piJwt,
			requestBody: &requests.ResetWeeklyDistributionRequest{
				WeekTotal: 10,
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestAddDevice(t *testing.T) {
	path := "/AddDevice"
	method := http.MethodPost

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           userJwt,
			requestBody: &requests.AddDeviceRequest{
				Name:       "smartTV",
				IpAddress:  "10.0.1.1",
				MacAddress: "00:00:00:00:00:00",
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
		{
			name:          "missing name",
			correctStatus: http.StatusBadRequest,
			jwt:           userJwt,
			requestBody: &requests.AddDeviceRequest{
				IpAddress:  "10.0.1.1",
				MacAddress: "00:00:00:00:00:00",
			},
		},
		{
			name:          "missing ip address",
			correctStatus: http.StatusBadRequest,
			jwt:           userJwt,
			requestBody: &requests.AddDeviceRequest{
				Name:       "smartTV",
				MacAddress: "00:00:00:00:00:00",
			},
		},
		{
			name:          "missing mac address",
			correctStatus: http.StatusBadRequest,
			jwt:           userJwt,
			requestBody: &requests.AddDeviceRequest{
				Name:      "smartTV",
				IpAddress: "10.0.1.1",
			},
		},
	}

	missingBodyTest(t, method, path)
	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}

func TestGetDevices(t *testing.T) {
	path := "/GetDevices"
	method := http.MethodGet

	testCases := []testCase{
		{
			name:          "valid request",
			correctStatus: http.StatusOK,
			jwt:           userJwt,
			expectedBody: &[]database.Device{
				{
					Id:            1,
					Name:          "smartTV",
					IpAddress:     "10.0.1.1",
					MacAddress:    "00:00:00:00:00:00",
					DateAdded:     time.Now().Format(time.DateOnly),
					IncidentCount: 0,
				},
			},
		},
		{
			name:          "missing jwt",
			correctStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid jwt",
			correctStatus: http.StatusUnauthorized,
			jwt:           "invalid.jwt.token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, buildTest(tc, method, path))
	}
}
