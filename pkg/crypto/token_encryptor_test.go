package crypto

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestNewTokenEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid 32-byte key",
			key:     []byte("12345678901234567890123456789012"),
			wantErr: false,
		},
		{
			name:    "valid short key (will be hashed)",
			key:     []byte("shortkey"),
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     []byte{},
			wantErr: true,
		},
		{
			name:    "nil key",
			key:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewTokenEncryptor(tt.key)
			if tt.wantErr {
				if err == nil {
					t.Error("NewTokenEncryptor() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("NewTokenEncryptor() unexpected error: %v", err)
				}
				if encryptor == nil {
					t.Error("NewTokenEncryptor() returned nil encryptor")
				}
				if len(encryptor.key) != 32 {
					t.Errorf("NewTokenEncryptor() key length = %d, want 32", len(encryptor.key))
				}
			}
		})
	}
}

func TestNewTokenEncryptorFromString(t *testing.T) {
	tests := []struct {
		name    string
		keyStr  string
		wantErr bool
	}{
		{
			name:    "valid string key",
			keyStr:  "my-secret-encryption-key",
			wantErr: false,
		},
		{
			name:    "empty string key",
			keyStr:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewTokenEncryptorFromString(tt.keyStr)
			if tt.wantErr {
				if err == nil {
					t.Error("NewTokenEncryptorFromString() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("NewTokenEncryptorFromString() unexpected error: %v", err)
				}
				if encryptor == nil {
					t.Error("NewTokenEncryptorFromString() returned nil encryptor")
				}
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("test-encryption-key-for-oauth-tokens")
	encryptor, err := NewTokenEncryptor(key)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple token",
			plaintext: "ya29.a0AfH6SMBx...",
		},
		{
			name:      "long token",
			plaintext: strings.Repeat("a", 1000),
		},
		{
			name:      "token with special characters",
			plaintext: "token!@#$%^&*()_+-=[]{}|;:',.<>?/~`",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode token",
			plaintext: "token-with-unicode-🔐-characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := encryptor.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error: %v", err)
			}

			// Empty plaintext should return empty ciphertext
			if tt.plaintext == "" && encrypted != "" {
				t.Error("Encrypt() empty plaintext should return empty ciphertext")
			}

			// Ciphertext should be different from plaintext (unless empty)
			if tt.plaintext != "" && encrypted == tt.plaintext {
				t.Error("Encrypt() ciphertext should differ from plaintext")
			}

			// Ciphertext should be base64 encoded
			if tt.plaintext != "" {
				_, err := base64.StdEncoding.DecodeString(encrypted)
				if err != nil {
					t.Errorf("Encrypt() ciphertext is not valid base64: %v", err)
				}
			}

			// Decrypt
			decrypted, err := encryptor.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error: %v", err)
			}

			// Decrypted should match original plaintext
			if decrypted != tt.plaintext {
				t.Errorf("Decrypt() = %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDecryptToken(t *testing.T) {
	encryptor, err := NewTokenEncryptorFromString("my-secret-key")
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	accessToken := "ya29.a0AfH6SMBxSampleAccessToken123456789"
	refreshToken := "1//0gSampleRefreshToken987654321"

	// Test access token
	encryptedAccess, err := encryptor.EncryptToken(accessToken)
	if err != nil {
		t.Fatalf("EncryptToken() error: %v", err)
	}

	decryptedAccess, err := encryptor.DecryptToken(encryptedAccess)
	if err != nil {
		t.Fatalf("DecryptToken() error: %v", err)
	}

	if decryptedAccess != accessToken {
		t.Errorf("DecryptToken() = %q, want %q", decryptedAccess, accessToken)
	}

	// Test refresh token
	encryptedRefresh, err := encryptor.EncryptRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("EncryptRefreshToken() error: %v", err)
	}

	decryptedRefresh, err := encryptor.DecryptRefreshToken(encryptedRefresh)
	if err != nil {
		t.Fatalf("DecryptRefreshToken() error: %v", err)
	}

	if decryptedRefresh != refreshToken {
		t.Errorf("DecryptRefreshToken() = %q, want %q", decryptedRefresh, refreshToken)
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	encryptor, err := NewTokenEncryptorFromString("my-secret-key")
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext string
		wantErr    error
	}{
		{
			name:       "invalid base64",
			ciphertext: "not-valid-base64!!!",
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "too short ciphertext",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("short")),
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "corrupted ciphertext",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("this is not encrypted data that can be decrypted")),
			wantErr:    ErrDecryptionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tt.ciphertext)
			if err == nil {
				t.Error("Decrypt() expected error, got nil")
			}
			if err != tt.wantErr {
				t.Errorf("Decrypt() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptionDeterminism(t *testing.T) {
	encryptor, err := NewTokenEncryptorFromString("my-secret-key")
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	plaintext := "test-token-123"

	// Encrypt the same plaintext twice
	encrypted1, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	encrypted2, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Ciphertexts should be different (due to random nonce)
	if encrypted1 == encrypted2 {
		t.Error("Encrypt() should produce different ciphertexts for same plaintext (non-deterministic)")
	}

	// But both should decrypt to the same plaintext
	decrypted1, err := encryptor.Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	decrypted2, err := encryptor.Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Errorf("Decrypt() failed to recover original plaintext")
	}
}

func TestHashToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "simple token",
			token: "simple-token-123",
		},
		{
			name:  "long token",
			token: strings.Repeat("a", 1000),
		},
		{
			name:  "empty token",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashToken(tt.token)

			// Hash should be base64 encoded
			_, err := base64.StdEncoding.DecodeString(hash)
			if err != nil {
				t.Errorf("HashToken() result is not valid base64: %v", err)
			}

			// Hash should be deterministic
			hash2 := HashToken(tt.token)
			if hash != hash2 {
				t.Error("HashToken() should produce same hash for same input")
			}

			// Different inputs should produce different hashes
			if tt.token != "" {
				differentHash := HashToken(tt.token + "different")
				if hash == differentHash {
					t.Error("HashToken() should produce different hashes for different inputs")
				}
			}
		})
	}
}

func TestGenerateRandomKey(t *testing.T) {
	key1, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey() error: %v", err)
	}

	if len(key1) != 32 {
		t.Errorf("GenerateRandomKey() key length = %d, want 32", len(key1))
	}

	key2, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey() error: %v", err)
	}

	// Keys should be different (extremely unlikely to be the same)
	if string(key1) == string(key2) {
		t.Error("GenerateRandomKey() should produce different keys")
	}
}

func TestGenerateRandomKeyBase64(t *testing.T) {
	keyStr, err := GenerateRandomKeyBase64()
	if err != nil {
		t.Fatalf("GenerateRandomKeyBase64() error: %v", err)
	}

	// Should be valid base64
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		t.Errorf("GenerateRandomKeyBase64() result is not valid base64: %v", err)
	}

	// Decoded key should be 32 bytes
	if len(key) != 32 {
		t.Errorf("GenerateRandomKeyBase64() decoded key length = %d, want 32", len(key))
	}
}

func TestDifferentKeysProduceDifferentResults(t *testing.T) {
	plaintext := "my-secret-token"

	encryptor1, _ := NewTokenEncryptorFromString("key1")
	encryptor2, _ := NewTokenEncryptorFromString("key2")

	encrypted1, err := encryptor1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	// Trying to decrypt with different key should fail
	_, err = encryptor2.Decrypt(encrypted1)
	if err == nil {
		t.Error("Decrypt() with different key should fail")
	}
}

func BenchmarkEncrypt(b *testing.B) {
	encryptor, _ := NewTokenEncryptorFromString("benchmark-key")
	plaintext := "ya29.a0AfH6SMBxSampleAccessToken123456789"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encryptor.Encrypt(plaintext)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	encryptor, _ := NewTokenEncryptorFromString("benchmark-key")
	plaintext := "ya29.a0AfH6SMBxSampleAccessToken123456789"
	encrypted, _ := encryptor.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encryptor.Decrypt(encrypted)
	}
}

func BenchmarkHashToken(b *testing.B) {
	token := "ya29.a0AfH6SMBxSampleAccessToken123456789"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HashToken(token)
	}
}
