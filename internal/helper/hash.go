package helper

import (
	"crypto/md5"
)

func GetMD5(data string) string {
	s := md5.Sum([]byte(data))
	return string(s[:])
}
