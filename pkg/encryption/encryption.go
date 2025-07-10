package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
)

func ToRSA(privateKeyPem, publicKeyPem string) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKeyByte := []byte(privateKeyPem)
	privateKeyBlock, _ := pem.Decode(privateKeyByte)
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	publicKeyByte := []byte(publicKeyPem)
	publicKeyBlock, _ := pem.Decode(publicKeyByte)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	return privateKey.(*rsa.PrivateKey), publicKey.(*rsa.PublicKey)
}

func ToAES(keyString string) []byte {
	key := []byte(keyString)
	if len(key) > 32 {
		log.Fatal("AES-256 key must be 32 bytes long")
	}
	return key
}

func EncryptAes(aesKey, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(plaintext)
	if err != nil {
		return nil, fmt.Errorf("cannot create aes-gcm cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot generate aes-gcm nonce: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("cannot generate aes-gcm nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func DecryptAes(aesKey, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("cannot create aes-gcm cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cannot generate aes-gcm nonce: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedMessage := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt: %w", err)
	}

	return plaintext, nil
}
