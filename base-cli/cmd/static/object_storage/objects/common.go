package objects

import "strings"

func ParseBucketNameAndObjectKey(path string) (bucketName string, objectKey string) {
	if path == "" {
		return "", ""
	}

	parts := strings.SplitN(path, "/", 2)

	bucketName = parts[0]

	if len(parts) == 2 {
		objectKey = parts[1]
	}

	return
}
