package handler

import (
	"math/rand"
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	codes = make(map[string]struct {
		Code       string
		Expiration time.Time
	})
)

func GenerateCode(email string) string {
	mutex.Lock()
	defer mutex.Unlock()

	code := generateRandomCode(6)
	codes[email] = struct {
		Code       string
		Expiration time.Time
	}{
		Code:       code,
		Expiration: time.Now().Add(5 * time.Minute), // Código válido por 5 minutos
	}

	return code
}

func VerifyCode(email, code string) bool {
	mutex.RLock()
	entry, exists := codes[email]
	mutex.RUnlock()

	if !exists || time.Now().After(entry.Expiration) {
		return false
	}

	return entry.Code == code
}


func RemoveCode(email string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(codes, email)
}

func generateRandomCode(n int) string {
	// Fuente de números aleatorios
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	const digits = "0123456789"
	code := make([]byte, n)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}
