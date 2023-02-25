package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"io/fs"

	"os"
)

/*
decrypts content of file <filename> with the secretKey

returns plaintext as []byte
*/
func EncryptFileSha256AES256(secretKey, plainText []byte, fileName string, fileMode fs.FileMode) error {

	cipherText, err := EncryptSha256AES256(secretKey, plainText)
	if err != nil {
		return err
	}

	// Writing ciphertext file
	return os.WriteFile(fileName, cipherText, fileMode)
}

/*
decrypts content of file <filename> with the secretKey

returns plaintext as []byte
*/
func DecryptFileSha256AES256(secretKey []byte, fileName string) ([]byte, error) {
	cipherText, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return DecryptSha256AES256(secretKey, cipherText)
}

/*
encodes a plainText with secretKey
the secretKey gets hashed with sha-256 sum to the AES-256 key of a fixed length of 32byte/256bit
*/
func EncryptSha256AES256(secretKey, plainText []byte) ([]byte, error) {
	// Create a new SHA-256 hash object
	hash := sha256.New()

	// Write the message to the hash object
	hash.Write(secretKey)

	// Get the SHA-256 hash as a byte slice
	hashBytes := hash.Sum(nil)

	// encrypt the plaintext with the hash and return
	return EncryptGcm(hashBytes, plainText)
}

/*
decode cipherText with AES-256 with the sha256 sum of the secretKey
*/
func DecryptSha256AES256(secretKey, cipherText []byte) ([]byte, error) {
	// Create a new SHA-256 hash object
	hash := sha256.New()

	// Write the message to the hash object
	hash.Write(secretKey)

	// Get the SHA-256 hash as a byte slice
	hashBytes := hash.Sum(nil)

	// decode the cipherText with the hash and return
	return DecryptGcm(hashBytes, cipherText)
}

/*
key, plainText
key must be len(16, 24, 32) choosing (AES-128, AES-192, AES-256)
*/
func EncryptGcm(key, plainText []byte) ([]byte, error) {
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generating random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Decrypt file
	return gcm.Seal(nonce, nonce, plainText, nil), nil

}

/*
key, cipherText
key must be len(16, 24, 32) choosing (AES-128, AES-192, AES-256)
*/
func DecryptGcm(key, cipherText []byte) ([]byte, error) {
	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	return gcm.Open(nil, nonce, cipherText, nil)
}
