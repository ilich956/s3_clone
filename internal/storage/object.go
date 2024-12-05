package storage

import (
	"errors"
	"io"
	"log"
	"os"
	"triple-s/internal/config"
	"triple-s/pkg/csvutil"
)

func PutObject(bucketName, objectName, contentType string, contentLength int64, objectData io.Reader) error {
	buckets := csvutil.ExtractBucketNamesCSV()

	// Check if bucket exists
	if _, ok := buckets[bucketName]; !ok {
		log.Println("Failed to put object: non-existent bucket")
		return errors.New("bucket does not exist")
	}

	if objectName == "objects.csv" {
		log.Println("Failed to put object: forbidden object name")
		return errors.New("forbidden object name")
	}

	// Update object metadata in CSV
	csvutil.AddObjectCSV(bucketName, objectName, contentType, contentLength)
	csvutil.UpdateBucketTimeCSV(bucketName)

	// Create and write to the file
	filePath := *config.Dir + "/" + bucketName + "/" + objectName
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Failed to create file:", err)
		return errors.New("failed to create file")
	}
	defer file.Close()

	if _, err := io.Copy(file, objectData); err != nil {
		log.Println("Failed to copy body content:", err)
		return errors.New("failed to write object data")
	}

	// Update bucket status to non-empty
	csvutil.UpdateBucketStatusCSV(bucketName)

	return nil
}

func GetObject(bucketName, objectName string) (io.ReadCloser, error) {
	buckets := csvutil.ExtractBucketNamesCSV()

	// Check if bucket exists
	if _, ok := buckets[bucketName]; !ok {
		log.Println("Failed to get object: non-existent bucket")
		return nil, errors.New("bucket does not exist")
	}

	// Open the object file for reading
	filePath := *config.Dir + "/" + bucketName + "/" + objectName
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed to open object:", err)
		return nil, errors.New("failed to open object")
	}

	return file, nil
}

func DeleteObject(bucketName, objectName string) error {
	buckets := csvutil.ExtractBucketNamesCSV()

	// Check if the bucket exists
	if _, ok := buckets[bucketName]; !ok {
		log.Println("Failed to delete object: non-existent bucket")
		return errors.New("bucket does not exist")
	}

	// Attempt to delete the object file
	objectPath := *config.Dir + "/" + bucketName + "/" + objectName
	if err := os.Remove(objectPath); err != nil {
		log.Println("Failed to delete object:", err)
		return errors.New("failed to delete object")
	}

	// Update CSV files to reflect the deletion
	csvutil.DeleteRowObjectCSV(bucketName, objectName)
	csvutil.UpdateBucketTimeCSV(bucketName)
	csvutil.UpdateBucketStatusCSV(bucketName)

	return nil
}
