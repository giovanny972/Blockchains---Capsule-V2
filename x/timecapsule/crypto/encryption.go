package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	
	"golang.org/x/crypto/scrypt"
)

// EncryptionManager handles encryption and decryption operations
type EncryptionManager struct {
	keySize int // AES key size (16, 24, or 32 bytes)
}

// NewEncryptionManager creates a new encryption manager with AES-256
func NewEncryptionManager() *EncryptionManager {
	return &EncryptionManager{
		keySize: 32, // AES-256
	}
}

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	Data      []byte `json:"data"`
	Nonce     []byte `json:"nonce"`
	Salt      []byte `json:"salt"`
	Algorithm string `json:"algorithm"`
}

// GenerateKey generates a cryptographically secure random key
func (em *EncryptionManager) GenerateKey() ([]byte, error) {
	key := make([]byte, em.keySize)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// DeriveKeyFromPassword derives an encryption key from a password using scrypt
func (em *EncryptionManager) DeriveKeyFromPassword(password string, salt []byte) ([]byte, error) {
	if len(salt) == 0 {
		salt = make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			return nil, fmt.Errorf("failed to generate salt: %w", err)
		}
	}
	
	key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, em.keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	
	return key, nil
}

// Encrypt encrypts data using AES-GCM
func (em *EncryptionManager) Encrypt(data []byte, key []byte) (*EncryptedData, error) {
	if len(key) != em.keySize {
		return nil, fmt.Errorf("invalid key size: expected %d, got %d", em.keySize, len(key))
	}
	
	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	
	return &EncryptedData{
		Data:      ciphertext,
		Nonce:     nonce,
		Algorithm: "AES-256-GCM",
	}, nil
}

// Decrypt decrypts data using AES-GCM
func (em *EncryptionManager) Decrypt(encData *EncryptedData, key []byte) ([]byte, error) {
	if len(key) != em.keySize {
		return nil, fmt.Errorf("invalid key size: expected %d, got %d", em.keySize, len(key))
	}
	
	if encData.Algorithm != "AES-256-GCM" {
		return nil, fmt.Errorf("unsupported algorithm: %s", encData.Algorithm)
	}
	
	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	
	// Decrypt data
	plaintext, err := gcm.Open(nil, encData.Nonce, encData.Data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	
	return plaintext, nil
}

// EncryptWithPassword encrypts data using a password-derived key
func (em *EncryptionManager) EncryptWithPassword(data []byte, password string) (*EncryptedData, error) {
	// Generate salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	
	// Derive key from password
	key, err := em.DeriveKeyFromPassword(password, salt)
	if err != nil {
		return nil, err
	}
	
	// Encrypt data
	encData, err := em.Encrypt(data, key)
	if err != nil {
		return nil, err
	}
	
	// Add salt to encrypted data
	encData.Salt = salt
	
	return encData, nil
}

// DecryptWithPassword decrypts data using a password-derived key
func (em *EncryptionManager) DecryptWithPassword(encData *EncryptedData, password string) ([]byte, error) {
	if len(encData.Salt) == 0 {
		return nil, fmt.Errorf("missing salt for password-based decryption")
	}
	
	// Derive key from password
	key, err := em.DeriveKeyFromPassword(password, encData.Salt)
	if err != nil {
		return nil, err
	}
	
	// Decrypt data
	return em.Decrypt(encData, key)
}

// HashData creates a SHA-256 hash of the data
func HashData(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// VerifyDataIntegrity verifies that the data matches the provided hash
func VerifyDataIntegrity(data []byte, expectedHash string) bool {
	actualHash := HashData(data)
	return actualHash == expectedHash
}

// SecureRandom generates cryptographically secure random bytes
func SecureRandom(size int) ([]byte, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// WipeKey securely wipes a key from memory
func WipeKey(key []byte) {
	for i := range key {
		key[i] = 0
	}
}