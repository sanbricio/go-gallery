package codeGeneratorHandler

import (
	"crypto/rand"
	"math/big"
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
var RandFunc = rand.Int

func GenerateCode(key string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	code, err := generateRandomCode(6)
	if err != nil {
		return "", err
	}
	codes[key] = struct {
		Code       string
		Expiration time.Time
	}{
		Code:       code,
		Expiration: NowFunc().Add(EXPIRATION_CODE), // Código válido por 5 minutos
	}

	return code, nil
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
	// Si ha expirado el código lo eliminamos y devolvemos false
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

func generateRandomCode(n int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, n)
	for i := range code {
		num, err := RandFunc(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}
	return string(code), nil
}
