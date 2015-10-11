package random

import (
	"crypto/rand"
)

func Generate(n int) string {
	//const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const alphanum = "2345789ABCDEFHJKMNPQRSTUVWXYZ"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
