package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, bits)
	PrintKeyPairAndAddress(privkey, &privkey.PublicKey)
	return privkey, &privkey.PublicKey
}

// SignTransaction creates a signature of the transaction using the private key
func SignTransaction(privkey *rsa.PrivateKey, transaction string) []byte {
	hashed := sha256.Sum256([]byte(transaction))
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privkey, crypto.SHA256, hashed[:])
	return signature
}

// VerifySignature checks the signature against the transaction using the public key
func VerifySignature(pubkey *rsa.PublicKey, transaction string, signature []byte) bool {
	hashed := sha256.Sum256([]byte(transaction))
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed[:], signature)
	return err == nil
}

// Address generates a unique address for a public key
func Address(pubkey *rsa.PublicKey) string {
	pubkeyBytes, _ := x509.MarshalPKIXPublicKey(pubkey)
	hashed := sha256.Sum256(pubkeyBytes)
	return hex.EncodeToString(hashed[:])
}

func PrintKeyPairAndAddress(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) {
	// Print the private key in PEM format
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	fmt.Println("Private Key:")
	fmt.Println(string(privKeyPEM))

	// Print the public key in PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	fmt.Println("Public Key:")
	fmt.Println(string(pubKeyPEM))

	// Print the address
	address := Address(pubKey)
	fmt.Println("Address:")
	fmt.Println(address)
}