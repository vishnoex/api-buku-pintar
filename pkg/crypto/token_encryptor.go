package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

var (
	// ErrInvalidKey is returned when the encryption key is invalid
	ErrInvalidKey = errors.New("invalid encryption key: key must be 32 bytes for AES-256")
	
	// ErrInvalidCiphertext is returned when the ciphertext is invalid or corrupted
	ErrInvalidCiphertext = errors.New("invalid ciphertext: data may be corrupted")
	
	// ErrEncryptionFailed is returned when encryption fails
	ErrEncryptionFailed = errors.New("encryption failed")
	
	// ErrDecryptionFailed is returned when decryption fails
	ErrDecryptionFailed = errors.New("decryption failed")
)

// TokenEncryptor provides methods to encrypt and decrypt OAuth tokens
// Uses AES-256-GCM for authenticated encryption
type TokenEncryptor struct {
	key []byte // 32 bytes for AES-256
}

// NewTokenEncryptor creates a new TokenEncryptor with the given key
// The key should be 32 bytes for AES-256
// If the key is a string, it will be hashed to ensure 32 bytes
func NewTokenEncryptor(key []byte) (*TokenEncryptor, error) {
	if len(key) == 0 {
		return nil, ErrInvalidKey
	}

	// If key is not 32 bytes, hash it to ensure proper length
	if len(key) != 32 {
		key = hashKeyTo32Bytes(key)
	}

	return &TokenEncryptor{
		key: key,
	}, nil
}

// NewTokenEncryptorFromString creates a TokenEncryptor from a string key
// The string will be hashed to produce a 32-byte key for AES-256
func NewTokenEncryptorFromString(keyString string) (*TokenEncryptor, error) {
	if keyString == "" {
		return nil, ErrInvalidKey
	}

	key := hashKeyTo32Bytes([]byte(keyString))
	return &TokenEncryptor{
		key: key,
	}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
// Returns base64-encoded ciphertext with nonce prepended
func (e *TokenEncryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil // Empty string doesn't need encryption
	}

	// Create AES cipher block
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrEncryptionFailed
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrEncryptionFailed
	}

	// Encrypt the plaintext
	// The nonce is prepended to the ciphertext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-256-GCM
// Returns the original plaintext
func (e *TokenEncryptor) Decrypt(ciphertextBase64 string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil // Empty string doesn't need decryption
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	// Create AES cipher block
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	// Check minimum ciphertext length
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// EncryptToken encrypts an OAuth access token
// This is a convenience method that wraps Encrypt
func (e *TokenEncryptor) EncryptToken(token string) (string, error) {
	return e.Encrypt(token)
}

// DecryptToken decrypts an OAuth access token
// This is a convenience method that wraps Decrypt
func (e *TokenEncryptor) DecryptToken(encryptedToken string) (string, error) {
	return e.Decrypt(encryptedToken)
}

// EncryptRefreshToken encrypts an OAuth refresh token
// This is a convenience method that wraps Encrypt
func (e *TokenEncryptor) EncryptRefreshToken(token string) (string, error) {
	return e.Encrypt(token)
}

// DecryptRefreshToken decrypts an OAuth refresh token
// This is a convenience method that wraps Decrypt
func (e *TokenEncryptor) DecryptRefreshToken(encryptedToken string) (string, error) {
	return e.Decrypt(encryptedToken)
}

// HashToken creates a SHA-256 hash of a token
// Used for token blacklisting (we don't store the actual token, just its hash)
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// hashKeyTo32Bytes hashes a key of any length to exactly 32 bytes for AES-256
func hashKeyTo32Bytes(key []byte) []byte {
	hash := sha256.Sum256(key)
	return hash[:]
}

// GenerateRandomKey generates a cryptographically secure random 32-byte key
// Useful for generating encryption keys
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateRandomKeyBase64 generates a random key and returns it as base64 string
func GenerateRandomKeyBase64() (string, error) {
	key, err := GenerateRandomKey()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
