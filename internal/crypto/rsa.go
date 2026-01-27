package crypto

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
)

const RSAKeySize = 2048

// GenerateRSAKeyPair generates a new RSA key pair
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
    }
    
    return privateKey, &privateKey.PublicKey, nil
}

// EncodePrivateKeyToPEM converts private key to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
    privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privateKeyBytes,
    })
    return string(privateKeyPEM)
}

// EncodePublicKeyToPEM converts public key to PEM format
func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return "", fmt.Errorf("failed to marshal public key: %w", err)
    }
    
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    return string(publicKeyPEM), nil
}

// DecodePrivateKeyFromPEM converts PEM to private key
func DecodePrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode([]byte(privateKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }
    
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }
    
    return privateKey, nil
}

// DecodePublicKeyFromPEM converts PEM to public key
func DecodePublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(publicKeyPEM))
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM block")
    }
    
    publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse public key: %w", err)
    }
    
    publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
    if !ok {
        return nil, fmt.Errorf("not an RSA public key")
    }
    
    return publicKey, nil
}

// EncryptWithPublicKey encrypts data with RSA public key (for AES key encryption)
func EncryptWithPublicKey(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
    ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
    if err != nil {
        return nil, fmt.Errorf("failed to encrypt with public key: %w", err)
    }
    return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with RSA private key
func DecryptWithPrivateKey(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt with private key: %w", err)
    }
    return plaintext, nil
}

// SignWithPrivateKey creates a digital signature
func SignWithPrivateKey(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
    // For digital signatures, we'll implement this in signature.go
    // This is a placeholder for now
    return nil, fmt.Errorf("use signature.go for signing")
}