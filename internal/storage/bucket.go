package storage

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"triple-s/internal/config"
	"triple-s/internal/utils/validation"
	"triple-s/pkg/csvutil"
)

type AllBucket struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
	BucketName       string `xml:"Name"`
	CreationDate     string `xml:"CreationDate"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

func PutBucket(bucketName string) error {
	// Check for existing bucket names
	buckets := csvutil.ExtractBucketNamesCSV()
	if _, exists := buckets[bucketName]; exists {
		return fmt.Errorf("failed to create new bucket: duplicate name")
	}

	// Validate bucket name
	if err := validation.ValidateBucketName(bucketName); err != nil {
		return fmt.Errorf("failed to create new bucket: %v", err)
	}

	// Create directory for the bucket
	if err := os.Mkdir(*config.Dir+"/"+bucketName, 0o755); err != nil {
		log.Println("Failed to create bucket directory: ", err)
		return fmt.Errorf("failed to create new bucket")
	}

	// Add bucket details to CSV
	csvutil.AddBucketCSV(bucketName)
	csvutil.CreateObjectCSV(bucketName)

	return nil
}

func GetBuckets() ([]byte, error) {
	buckets := csvutil.ExtractBucketNamesCSV()

	allBuckets := AllBucket{}
	for bucketName, bucketData := range buckets {
		newBucket := Bucket{
			BucketName:       bucketName,
			CreationDate:     bucketData[0],
			LastModifiedTime: bucketData[1],
			Status:           bucketData[2],
		}
		allBuckets.Buckets = append(allBuckets.Buckets, newBucket)
	}

	// Marshal the buckets into XML format
	output, err := xml.MarshalIndent(allBuckets, "  ", "    ")
	if err != nil {
		log.Print("Failed to marshal buckets: ", err)
		return nil, err
	}
	return output, nil
}

func DeleteBucket(bucketName string) error {
	buckets := csvutil.ExtractBucketNamesCSV()
	bucketData, exists := buckets[bucketName]

	// Delete bucket folder
	// err := os.RemoveAll(*config.Dir + "/" + bucketName)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete bucket folder: %w", err)
	// }

	// Check if bucket exists
	if !exists {
		log.Println("Failed to delete bucket: non-existent bucket")
		return errors.New("bucket does not exist")
	}

	// Check if bucket is empty
	if bucketData[2] == "non-empty" {
		log.Println("Failed to delete bucket: bucket is not empty")
		return errors.New("bucket is not empty")
	}

	// Check if it's a directory
	info, err := os.Stat(*config.Dir + "/" + bucketName)
	if err != nil {
		log.Println("IS DIR")
		if os.IsNotExist(err) {
			return fmt.Errorf("No such file or directory: %s\n", bucketName)
		}
		return fmt.Errorf("Error checking file info: %v\n", err)
	}

	// If it's a directory, remove it
	if info.IsDir() {
		log.Println("IS DIR")
		err := os.RemoveAll(*config.Dir + "/" + bucketName) // RemoveAll will delete the directory and its contents
		if err != nil {
			log.Println("Error removing directory: %v\n", err)
			return fmt.Errorf("Error removing directory: %v\n", err)
		}
		log.Println("Successfully deleted directory: %s\n", bucketName)
	}

	// Delete the bucket from CSV
	if err := csvutil.DeleteRowBucketCSV(bucketName); err != nil {
		log.Println("Failed to delete bucket:", err)
		return errors.New("failed to delete bucket")
	}

	return nil
}
