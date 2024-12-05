package validation

import (
	"fmt"
	"regexp"
)

func ValidateBucketName(bucketName string) error {
	// Check length of bucket name
	if len(bucketName) < 3 || len(bucketName) > 63 {
		return fmt.Errorf("bucket names should be between 3 and 63 characters long")
	}

	pattern := `^[a-z0-9]+([-.]?[a-z0-9]+)*$`
	ipPattern := `^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`
	patternRegex := regexp.MustCompile(pattern)
	ipPatternRegex := regexp.MustCompile(ipPattern)

	if ipPatternRegex.MatchString(bucketName) {
		return fmt.Errorf("bucket name must not be formatted as an IP address (e.g., 192.168.0.1)")
	}

	if bucketName[0] == '.' || bucketName[0] == '-' || bucketName[len(bucketName)-1] == '.' || bucketName[len(bucketName)-1] == '-' {
		return fmt.Errorf("bucket name must not begin or end with a hyphen or period")
	}

	if !patternRegex.MatchString(bucketName) {
		return fmt.Errorf("bucket name contains invalid characters or consecutive periods/dashes")
	}

	return nil
}
