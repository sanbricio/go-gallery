package codeGeneratorRepository

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync"
	"time"
)

const (
	CodeGeneratorMemoryRepositoryKey string = "CodeGeneratorMemoryRepository"
	DEFAULT_EXPIRATION_CODE          int    = 5
	DEFAULT_CLEANUP_INTERVAL         int    = 1
)

type CodeGeneratorMemoryRepository struct {
	expirationCode  time.Duration
	cleanupInterval time.Duration
}

func NewCodeGeneratorMemory(args map[string]string) *CodeGeneratorMemoryRepository {
	expirationCode, err := strconv.Atoi(args["CODE_GENERATOR_EXPIRATION_CODE"])
	if err != nil {
		expirationCode = DEFAULT_EXPIRATION_CODE
	}

	cleanupInterval, err := strconv.Atoi(args["CODE_GENERATOR_CLEANUP_INTERVAL"])
	if err != nil {
		cleanupInterval = DEFAULT_CLEANUP_INTERVAL
	}

	codeGenerator := &CodeGeneratorMemoryRepository{
		expirationCode:  time.Duration(expirationCode) * time.Minute,
		cleanupInterval: time.Duration(cleanupInterval) * time.Minute,
	}

	codeGenerator.StartAutoCleanup()

	return codeGenerator
}

var (
	mutex sync.RWMutex
	codes = make(map[string]struct {
		Code       string
		Expiration time.Time
	})

	NowFunc  = time.Now
	RandFunc = rand.Int
)

// This implementation could be a redis in the future
func (c *CodeGeneratorMemoryRepository) GenerateCode(key string) (string, error) {
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
		Expiration: NowFunc().Add(c.expirationCode),
	}

	return code, nil
}

func (c *CodeGeneratorMemoryRepository) VerifyCode(key, code string) bool {
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
		c.removeCode(key)
		return false
	}

	return entry.Code == code
}

func (c *CodeGeneratorMemoryRepository) removeCode(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(codes, key)
}

func (c *CodeGeneratorMemoryRepository) StartAutoCleanup() {
	go func() {
		for {
			time.Sleep(c.cleanupInterval)
			cleanupExpiredCodes()
		}
	}()
}

func cleanupExpiredCodes() {
	mutex.Lock()
	defer mutex.Unlock()
	now := NowFunc()
	for key, entry := range codes {
		if now.After(entry.Expiration) {
			delete(codes, key)
		}
	}
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
