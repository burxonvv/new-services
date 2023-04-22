package etc

import (
	"crypto/rand"
	"io"
	"log"
)

var (
	// table for code generator
	table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

// GenerateCode is function generating n-digit random code
func GenerateCode(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		log.Println("failed to read max: ", err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return string(b)
}
