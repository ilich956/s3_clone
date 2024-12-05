package handler

import (
	"io"
	"log"
	"net/http"
	"strings"
	"triple-s/internal/storage"
	"triple-s/internal/utils/response"
)

func HandlePutBucket(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if strings.Count(r.URL.Path, "/") != 1 {
		response.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract bucket name
	bucketName := r.PathValue("BucketName")

	// Attempt to create bucket
	if err := storage.PutBucket(bucketName); err != nil {
		log.Println(err)
		response.SendError(w, http.StatusConflict, err.Error())
		return
	}

	// Send success message
	response.SendSuccess(w, http.StatusOK, "Bucket successfully created")
	log.Print("Bucket successfully created")
}

func HandleGetBucket(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if r.URL.Path != "/" || strings.Count(r.URL.Path, "/") != 1 {
		response.SendError(w, http.StatusMethodNotAllowed, "Wrong URL")
		return
	}

	// Retrieve and marshal bucket data
	output, err := storage.GetBuckets()
	if err != nil {
		log.Print("Failed to retrieve buckets: ", err)
		response.SendError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	// Send successful XML response
	w.Header().Set("Content-Type", "application/xml")
	w.Write(output)
}

func HandleDeleteBucket(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if strings.Count(r.URL.Path, "/") != 1 {
		response.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract bucket name from URL path
	bucketName := r.PathValue("BucketName")

	// Attempt to delete the bucket
	err := storage.DeleteBucket(bucketName)
	switch {
	case err == nil:
		// Bucket successfully deleted
		w.WriteHeader(http.StatusNoContent)
	case err.Error() == "bucket does not exist":
		// Bucket not found
		response.SendError(w, http.StatusNotFound, "Failed to delete bucket: non-existent bucket")
	case err.Error() == "bucket is not empty":
		// Bucket is not empty
		response.SendError(w, http.StatusConflict, "Failed to delete bucket: bucket is not empty")
	default:
		// Internal server error
		response.SendError(w, http.StatusInternalServerError, "Failed to delete bucket")
	}
}

func HandlePutObject(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if strings.Count(r.URL.Path, "/") != 2 {
		response.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract bucket and object names from URL path
	bucketName := r.PathValue("BucketName")
	objectName := r.PathValue("ObjectKey")

	// Get content type and length from headers
	contentType := r.Header.Get("Content-Type")
	contentLength := r.ContentLength

	// Call storage.PutObject to handle object creation
	err := storage.PutObject(bucketName, objectName, contentType, contentLength, r.Body)
	switch {
	case err == nil:
		// Object successfully uploaded
		w.WriteHeader(http.StatusOK)
		response.SendSuccess(w, http.StatusOK, "Object successfully uploaded")
	case err.Error() == "bucket does not exist":
		// Bucket not found
		response.SendError(w, http.StatusNotFound, "Failed to put object: non-existent bucket")
	case err.Error() == "forbidden object name":
		response.SendError(w, http.StatusNotFound, "Failed to put object: forbidden object name")
	case err.Error() == "failed to create file" || err.Error() == "failed to write object data":
		response.SendError(w, http.StatusInternalServerError, "Internal server error")
	default:
		// Other errors
		response.SendError(w, http.StatusInternalServerError, "Failed to upload object")
	}
}

func HandleGetObject(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if strings.Count(r.URL.Path, "/") != 2 {
		response.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract bucket and object names from URL path
	bucketName := r.PathValue("BucketName")
	objectName := r.PathValue("ObjectKey")

	// Call storage.GetObject to retrieve the object
	objectData, err := storage.GetObject(bucketName, objectName)
	if err != nil {
		switch err.Error() {
		case "bucket does not exist":
			// Bucket not found
			response.SendError(w, http.StatusNotFound, "Failed to get object: non-existent bucket")
		case "failed to open object":
			// Object not found or cannot be opened
			response.SendError(w, http.StatusInternalServerError, "Failed to get object")
		default:
			// Other errors
			response.SendError(w, http.StatusInternalServerError, "Failed to retrieve object")
		}
		return
	}
	defer objectData.Close()

	// Stream the object data to the response
	if _, err := io.Copy(w, objectData); err != nil {
		log.Println("Failed to copy object to response:", err)
		response.SendError(w, http.StatusInternalServerError, "Internal error")
	}
}

// CHECK IF OBJECT EXISTS
func HandleDeleteObject(w http.ResponseWriter, r *http.Request) {
	// Validate URL structure
	if strings.Count(r.URL.Path, "/") != 2 {
		response.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract bucket and object names from URL path
	bucketName := r.PathValue("BucketName")
	objectName := r.PathValue("ObjectKey")

	// Call storage.DeleteObject to delete the object
	err := storage.DeleteObject(bucketName, objectName)
	if err != nil {
		switch err.Error() {
		case "bucket does not exist":
			// Bucket not found
			response.SendError(w, http.StatusNotFound, "Failed to delete object: non-existent bucket")
		case "failed to delete object":
			// Object not found or cannot be deleted
			response.SendError(w, http.StatusNotFound, "Failed to delete object")
		default:
			// Other errors
			response.SendError(w, http.StatusInternalServerError, "Failed to delete object")
		}
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}
