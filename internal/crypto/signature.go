package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// CreateSignature signs data with private key
func CreateSignature(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Hash the data
	hash := sha256.Sum256(data)

	// Sign the hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create signature: %w", err)
	}

	return signature, nil
}

// VerifySignature verifies a signature with public key
func VerifySignature(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	// Hash the data
	hash := sha256.Sum256(data)

	// Verify signature
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
