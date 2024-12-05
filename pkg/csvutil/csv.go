package csvutil

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"triple-s/internal/config"
)

func CreateNewBucketCSV() {
	csvFile, err := os.Create(*config.Dir + "/buckets.csv")
	if err != nil {
		log.Fatal("Failed to create buckets.csv: ", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headers := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
	writer.Write(headers)
}

func AddBucketCSV(bucketName string) {
	csvFile, err := os.OpenFile(*config.Dir+"/buckets.csv", os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatal("Failed to create buckets.csv: ", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()
	currentTime := time.Now().Format(time.RFC3339)
	row := []string{bucketName, currentTime, currentTime, "empty"} // add checking if there is any files in bucket,if no put empty, otherwise non-empty
	writer.Write(row)
}

func UpdateBucketStatusCSV(bucketName string) error {
	csvFilePath := *config.Dir + "/buckets.csv"

	// Read existing CSV data
	var rows [][]string

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err = reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %v", err)
	}

	// Check if the object key already exists and update the row if it does
	for i, row := range rows {
		if row[0] == bucketName {
			file, err := os.Open(*config.Dir + "/" + bucketName)
			if err != nil {
				log.Fatal("Failed to open bucket dir", err)
			}

			fileCount, err := file.Readdirnames(-1)
			if err != nil {
				log.Fatal("Failed to read bucket dir", err)
			}

			status := "empty"
			if len(fileCount) > 1 {
				status = "non-empty"
			}
			rows[i][3] = status
			break
		}
	}

	// open and  truncate contents
	csvFile, err = os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file for writing: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// write all rows back to the CSV file
	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write to CSV: %v", err)
	}

	return nil
}

func DeleteRowBucketCSV(bucketName string) error {
	csvFilePath := *config.Dir + "/buckets.csv"

	// Open the CSV file for reading and writing
	csvFile, err := os.OpenFile(csvFilePath, os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open buckets.csv: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read buckets.csv: %w", err)
	}

	// Move the file pointer to the beginning of the file
	csvFile.Seek(0, 0)

	writer := csv.NewWriter(csvFile)

	// Write records back, skipping the one to delete
	for _, record := range records {
		if record[0] == bucketName {
			continue // Skip the record to delete
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV file: %w", err)
		}
	}

	// Truncate the file to remove old content after writing
	// Move to the current position after writing
	offset, err := csvFile.Seek(0, 1)
	if err != nil {
		return fmt.Errorf("failed to seek to current position: %w", err)
	}

	// Truncate the file from the current offset
	if err := csvFile.Truncate(offset); err != nil {
		return fmt.Errorf("failed to truncate CSV file: %w", err)
	}

	writer.Flush() // Ensure all buffered operations are applied
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil // Return nil if successful
}

func UpdateBucketTimeCSV(bucketName string) error {
	csvFilePath := *config.Dir + "/buckets.csv"

	// Read existing CSV data
	var rows [][]string

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err = reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %v", err)
	}

	// Check if the object key already exists and update the row if it does
	for i, row := range rows {
		if row[0] == bucketName { // Assuming objectKey is in the first column
			// Update the existing row
			currentTime := time.Now().Format(time.RFC3339)
			rows[i][2] = currentTime
			break
		}
	}

	// open and  truncate contents
	csvFile, err = os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file for writing: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write all rows back to the CSV file
	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write to CSV: %v", err)
	}

	return nil
}

func ExtractBucketNamesCSV() map[string][]string {
	csvFile, err := os.Open(*config.Dir + "/buckets.csv")
	if err != nil {
		log.Fatal("Failed to open buckets.csv: ", err)
	}
	defer csvFile.Close()

	data, err := io.ReadAll(csvFile)
	if err != nil {
		log.Fatal("Failed to read buckets.csv: ", err)
	}

	reader := csv.NewReader(bytes.NewReader(data))

	buckets := make(map[string][]string)
	// var bucketNames []string

	// Skip first row
	if _, err := reader.Read(); err != nil {
		log.Fatal("Failed to read CSV data")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Failed to read CSV data")
		}
		creationTime := record[1]
		lastModifiedTime := record[2]
		status := record[3]
		buckets[record[0]] = append(buckets[record[0]], creationTime, lastModifiedTime, status)
		// bucketNames = append(bucketNames, record[0])
	}
	// return bucketNames[1:]
	return buckets
}

func AddObjectCSV(bucketName, objectKey, contentType string, size int64) error {
	csvFilePath := *config.Dir + "/" + bucketName + "/objects.csv"

	// Read existing CSV data
	var rows [][]string
	exists := false

	csvFile, err := os.Open(csvFilePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	rows, err = reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV data: %v", err)
	}

	// Check if the object key already exists and update the row if it does
	for i, row := range rows {
		if i == 0 {
			continue // skip first row
		}
		if row[0] == objectKey {
			// Update the existing row
			currentTime := time.Now().Format(time.RFC3339)
			rows[i] = []string{objectKey, strconv.FormatInt(size, 10), contentType, currentTime}
			exists = true
			break
		}
	}

	// If the object key does not exist, append a new row
	if !exists {
		currentTime := time.Now().Format(time.RFC3339)
		rows = append(rows, []string{objectKey, strconv.FormatInt(size, 10), contentType, currentTime})
	}

	/// open and  truncate contents
	csvFile, err = os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file for writing: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write all rows back to the CSV file
	if err := writer.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write to CSV: %v", err)
	}

	return nil
}

func CreateObjectCSV(bucketName string) {
	csvFile, err := os.Create(*config.Dir + "/" + bucketName + "/objects.csv")
	if err != nil {
		log.Fatal("Failed to create objects.csv: ", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	headers := []string{"ObjectKey", "Size", "ContentType", "LastModified"}
	writer.Write(headers)
}

func DeleteRowObjectCSV(bucketName, objectName string) error {
	csvFilePath := *config.Dir + "/" + bucketName + "/objects.csv"

	// Open the CSV file for reading and writing
	csvFile, err := os.OpenFile(csvFilePath, os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open objects.csv: %w", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read objects.csv: %w", err)
	}

	// Move file pointer to the beginning of the file
	csvFile.Seek(0, 0)

	writer := csv.NewWriter(csvFile) // Write first row of csv
	if err := writer.Write(records[0]); err != nil {
		return fmt.Errorf("failed to write record to CSV file: %w", err)
	}

	records = records[1:] // Skip first row

	// Write records
	for _, record := range records {
		if record[0] == objectName {
			continue // Skip the record to delete
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV file: %w", err)
		}
	}

	// Move to the current position after writing
	offset, err := csvFile.Seek(0, 1)
	if err != nil {
		return fmt.Errorf("failed to seek to current position: %w", err)
	}

	// Truncate the file from the current offset
	if err := csvFile.Truncate(offset); err != nil {
		return fmt.Errorf("failed to truncate CSV file: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}
