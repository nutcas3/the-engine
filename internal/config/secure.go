package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// SecureEnv manages secure environment variables
type SecureEnv struct {
	masterKey []byte
}

// NewSecureEnv creates a new secure environment manager
func NewSecureEnv() (*SecureEnv, error) {
	key := os.Getenv("ENGINE_MASTER_KEY")
	if key == "" {
		return nil, fmt.Errorf("ENGINE_MASTER_KEY environment variable is required")
	}

	// Derive a proper key from the master key
	hash := sha256.Sum256([]byte(key))

	return &SecureEnv{
		masterKey: hash[:],
	}, nil
}

// Encrypt encrypts a value using AES-GCM
func (s *SecureEnv) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a value using AES-GCM
func (s *SecureEnv) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, ciphertextBytes, nonce, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GetSecureEnv retrieves and decrypts an environment variable
func (s *SecureEnv) GetSecureEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not found", key)
	}

	// Check if the value is encrypted (starts with "enc:")
	if after, ok := strings.CutPrefix(value, "enc:"); ok {
		return s.Decrypt(after)
	}

	return value, nil
}

// SecureEnvFile represents a secure environment file
type SecureEnvFile struct {
	Variables map[string]string `json:"variables"`
}

// LoadSecureEnvFile loads and decrypts environment variables from a file
func LoadSecureEnvFile(path string) error {
	secure, err := NewSecureEnv()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var envFile SecureEnvFile
	if err := json.Unmarshal(data, &envFile); err != nil {
		return err
	}

	for key, value := range envFile.Variables {
		decrypted, err := secure.Decrypt(value)
		if err != nil {
			return fmt.Errorf("failed to decrypt %s: %w", key, err)
		}
		os.Setenv(key, decrypted)
	}

	return nil
}

// GenerateMasterKey generates a new master key for encryption
func GenerateMasterKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
