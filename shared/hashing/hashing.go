package hashing

import (
	"crypto/md5"
	"encoding/hex"
)

// Source: http://stackoverflow.com/questions/2377881/how-to-get-a-md5-hash-from-a-string-in-golang

// MD5Hash returns an MD5
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
