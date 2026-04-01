package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	secret := "my-super-secret-key-123!"
	plaintext := "192.168.1.100"

	// Test Encryption
	encrypted, err := Encrypt(plaintext, secret)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}
	if encrypted == "" {
		t.Fatal("Encrypted string is empty")
	}
	if encrypted == plaintext {
		t.Fatal("Encrypted string matches plaintext")
	}

	// Test Decryption
	decrypted, err := Decrypt(encrypted, secret)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected decrypted text to be %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongSecret(t *testing.T) {
	secret := "my-super-secret-key-123!"
	wrongSecret := "wrong-secret"
	plaintext := "192.168.1.100"

	encrypted, err := Encrypt(plaintext, secret)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	_, err = Decrypt(encrypted, wrongSecret)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong secret, got nil")
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	secret := "my-secret"
	
	_, err := Decrypt("invalid-base64-string", secret)
	if err == nil {
		t.Fatal("Expected error when decrypting invalid base64, got nil")
	}

	_, err = Decrypt("YQ==", secret) // Valid base64, but too short to contain nonce
	if err == nil {
		t.Fatal("Expected error when decrypting too short ciphertext, got nil")
	}
}
