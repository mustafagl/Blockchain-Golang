package blockchain

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, bits)
	//PrintKeyPairAndAddress(privkey, &privkey.PublicKey)
	return privkey, &privkey.PublicKey
}

// SignTransaction creates a signature of the transaction using the private key
func SignTransaction(privkey *rsa.PrivateKey, transaction []byte) []byte {
	hashed := sha256.Sum256(transaction)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privkey, crypto.SHA256, hashed[:])
	return signature
}

// VerifySignature checks the signature against the transaction using the public key
func VerifySignature(pubkey *rsa.PublicKey, transaction []byte, signature []byte) bool {
	hashed := sha256.Sum256(transaction)
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

func SaveKeyPairAndAddressToJSON(filename string, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, address string) error {
	// Convert the private key to PEM format
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	// Convert the public key to PEM format
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("Error marshalling public key: %v", err)
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	})

	// Create a map to hold the key pair and address
	keyPairAndAddressMap := map[string]string{
		"private_key": string(privKeyPEM),
		"public_key":  string(pubKeyPEM),
		"address":     address,
	}

	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(keyPairAndAddressMap, "", "    ")
	if err != nil {
		return fmt.Errorf("Error marshalling JSON: %v", err)
	}

	// Write the JSON data to a file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing JSON file: %v", err)
	}

	return nil
}

func CreateAndSaveNewWallet() (*rsa.PrivateKey, *rsa.PublicKey, string) {
	private_key, public_key := GenerateKeyPair(2048)
	address := Address(public_key)
	err := SaveKeyPairAndAddressToJSON("keypair_and_address.json", private_key, public_key, address)
	if err != nil {
		fmt.Println("Error saving to JSON:", err)
	}
	return private_key, public_key, address
}
