package codeGeneratorHandler

import (
	"math/rand"
	"sync"
	"time"
)

const (
	EXPIRATION_CODE time.Duration = 5 * time.Minute
)

var (
	mutex sync.RWMutex
	codes = make(map[string]struct {
		Code       string
		Expiration time.Time
	})
)

var NowFunc = time.Now

func GenerateCode(key string) string {
	mutex.Lock()
	defer mutex.Unlock()

	code := generateRandomCode(6)
	codes[key] = struct {
		Code       string
		Expiration time.Time
	}{
		Code:       code,
		Expiration: NowFunc().Add(EXPIRATION_CODE), // Código válido por 5 minutos
	}

	return code
}

func VerifyCode(key, code string) bool {
	if key == "" || code == "" {
		return false
	}

	mutex.RLock()
	entry, exists := codes[key]
	mutex.RUnlock()

	if !exists {
		return false
	}
	// Si ha expirado el codigo lo eliminamos y devolvemos false
	if NowFunc().After(entry.Expiration) {
		RemoveCode(key)
		return false
	}

	return entry.Code == code
}

func RemoveCode(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(codes, key)
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
