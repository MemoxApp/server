package local

import "math/rand"

// GenerateRandomString 生成随机字符串
func GenerateRandomString(length int) string {
	seed := "abcdef0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = seed[rand.Intn(len(seed))]
	}
	return string(b)
}
