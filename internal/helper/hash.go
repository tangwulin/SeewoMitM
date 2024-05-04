package helper

import "crypto/md5"

func MD5Sum(data []byte) string {
	original := md5.Sum(data)
	return string(original[:])
}
