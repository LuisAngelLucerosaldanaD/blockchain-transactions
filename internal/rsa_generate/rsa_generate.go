package rsa_generate

import (
	"blockchain-transactions/internal/logger"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

func Execute() (string, string, error) {
	strPrivate := ""
	strPublic := ""
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		logger.Error.Printf("couldn't generate rsaPrivate: %v", err)
		return strPrivate, strPublic, err
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		logger.Error.Printf("couldn't generate rsaPublic: %v", err)
		return strPrivate, strPublic, err
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)
	strPrivate = string(privateKeyBytes)
	strPublic = string(publicKeyBytes)
	return strPrivate, strPublic, nil
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	return pem.EncodeToMemory(&privBlock)

}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}
