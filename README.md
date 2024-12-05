# AWS S3 Clone


## Overview

AWS S3 Clone - (Simplified Storage Service) is a simplified implementation of the Amazon S3 storage service. It demonstrates core cloud storage principles by providing a RESTful API for bucket management and object operations. Designed as an educational tool, the project introduces basic networking concepts, HTTP server setup, and REST API design using Go.

# API Endpoints Documentation

## Bucket Operations

### 1. Create a Bucket
- **HTTP Method:** `PUT`
- **Endpoint:** `/{BucketName}`
- **Request Body:** Empty  
- **Behavior:**
  - Validate the bucket name against [S3 naming conventions](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html).
  - Ensure bucket name is unique.
  - If valid and unique, create a new bucket entry in the metadata.
  - **Response:**
    - `200 OK`: Bucket created successfully.
    - `400 Bad Request`: Invalid bucket name.
    - `409 Conflict`: Bucket already exists.

---

### 2. List All Buckets
- **HTTP Method:** `GET`
- **Endpoint:** `/`
- **Behavior:**
  - Fetch all bucket metadata from storage.
  - Return an [XML response](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html#API_ListBuckets_Examples) listing all buckets.
  - **Response:**
    - `200 OK`: List of all buckets.

---

### 3. Delete a Bucket
- **HTTP Method:** `DELETE`
- **Endpoint:** `/{BucketName}`
- **Behavior:**
  - Verify if the bucket exists and is empty.
  - Delete the bucket from metadata storage.
  - **Response:**
    - `204 No Content`: Bucket deleted successfully.
    - `404 Not Found`: Bucket does not exist.
    - `409 Conflict`: Bucket is not empty.

---

## Object Operations

### 1. Upload a New Object
- **HTTP Method:** `PUT`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Request Body:** Binary file content.
- **Headers:**  
  - `Content-Type`: MIME type of the object.  
  - `Content-Length`: Size of the object in bytes.  
- **Behavior:**
  - Validate the bucket and object key.
  - Save the file to `data/{BucketName}/{ObjectKey}`.
  - Update metadata in `data/{BucketName}/objects.csv`.
  - Overwrite if the object already exists.
  - **Response:**
    - `200 OK`: Object uploaded successfully.
    - `400 Bad Request`: Invalid input.
    - `404 Not Found`: Bucket does not exist.

---

### 2. Retrieve an Object
- **HTTP Method:** `GET`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Behavior:**
  - Validate the bucket and object existence.
  - Return the object's binary content.
  - **Response:**
    - `200 OK`: Object retrieved successfully.
    - `404 Not Found`: Object or bucket does not exist.

---

### 3. Delete an Object
- **HTTP Method:** `DELETE`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Behavior:**
  - Validate the bucket and object existence.
  - Remove the object from storage and metadata.
  - **Response:**
    - `204 No Content`: Object deleted successfully.
    - `404 Not Found`: Object or bucket does not exist.

---

## General Notes
- **Bucket Naming Rules:**
  - 3-63 characters.
  - Lowercase letters, numbers, hyphens, and periods allowed.
  - No IP address format (e.g., `192.168.1.1`).
  - No consecutive periods or dashes.
- **Validation Implementation:**
  - Regular expressions enforce naming conventions.
  - Uniqueness checked against stored metadata.

Refer to [S3 API Documentation](https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html) for further details and examples.
