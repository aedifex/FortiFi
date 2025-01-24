package handler

import (
	"log"
	"net/http"
)

func NotifyIntrusionHandler(writer http.ResponseWriter, request *http.Request){
	log.Printf("Request from %s", request.RemoteAddr)
	if request.Method != "POST" {
		log.Printf("Invalid Method from %s", request.RemoteAddr)
		writer.WriteHeader(405)
		response := "Method Not Allowed"
		writer.Write([]byte(response))
		return
	}
	// Correct flow:
	// Request body should include information about the intrusion and device uid as it is entered in database 
	// Server should update database with event info
	// Server should send notification to user accordingly
	//    Need to get user associated with specific pi from database
}