package response

import (
	"encoding/xml"
	"log"
	"net/http"
)

type Error struct {
	XMLName          xml.Name `xml:"Error"`
	ErrorCode        int      `xml:"ErrorCode"`
	ErrorDescription string   `xml:"ErrorDescription"`
}

type Success struct {
	XMLName xml.Name `xml:"Response"`
	Message string   `xml:"Message"`
}

func SendError(w http.ResponseWriter, statusCode int, errorDescription string) {
	error := Error{
		ErrorCode:        statusCode,
		ErrorDescription: errorDescription,
	}

	output, err := xml.MarshalIndent(error, "  ", "    ")
	if err != nil {
		log.Print("Failed to marshal error ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	w.Write(output)
}

func SendSuccess(w http.ResponseWriter, statusCode int, message string) {
	successResponse := Success{
		Message: message,
	}

	output, err := xml.MarshalIndent(successResponse, "  ", "    ")
	if err != nil {
		log.Print("Failed to marshal success response: ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(output)
}
