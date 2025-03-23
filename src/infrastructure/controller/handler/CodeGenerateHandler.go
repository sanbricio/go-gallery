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

var NowFunc = time.Now

func GenerateCode(email string) string {
	mutex.Lock()
	defer mutex.Unlock()

	code := generateRandomCode(6)
	codes[email] = struct {
		Code       string
		Expiration time.Time
	}{
		Code:       code,
		Expiration: NowFunc().Add(5 * time.Minute), // Código válido por 5 minutos
	}

	return code
}

func VerifyCode(email, code string) bool {
	if code == "" || email == "" {
		return false
	}

	mutex.RLock()
	entry, exists := codes[email]
	mutex.RUnlock()

	if !exists {
		return false
	}
	// Si ha expirado el codigo lo eliminamos y devolvemos false
	if NowFunc().After(entry.Expiration) {
		RemoveCode(email)
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
	source := rand.NewSource(NowFunc().UnixNano())
	r := rand.New(source)

	const digits = "0123456789"
	code := make([]byte, n)
	for i := range code {
		code[i] = digits[r.Intn(len(digits))]
	}
	return string(code)
}
