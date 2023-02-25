package encryption_test

import (
	"testing"

	"github.com/ricoschulte/go-myapps/encryption"
	"gotest.tools/assert"
)

func TestEncryptDecryptGcm(t *testing.T) {
	tests := []struct {
		name      string
		keyLength int

		plaintext string
		secretKey string
	}{
		{
			name:      "AES-128",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey22222",
			keyLength: 16,
		},
		{
			name:      "AES-192",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey2222212345648",
			keyLength: 24,
		},
		{
			name:      "AES-256",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey22222mysecretkey22222",
			keyLength: 32,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cipher, _ := encryption.EncryptGcm([]byte(test.secretKey), []byte(test.plaintext))
			plain, _ := encryption.DecryptGcm([]byte(test.secretKey), cipher)
			assert.Equal(t, test.plaintext, string(plain))
		})
	}
}

func TestDecryptGcm(t *testing.T) {
	tests := []struct {
		name      string
		keyLength int
		cipher    []byte
		plaintext string
		secretKey string
	}{
		{
			name:      "AES-128",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey22222",
			cipher:    []byte{54, 183, 205, 45, 205, 39, 173, 202, 115, 7, 93, 241, 124, 36, 161, 179, 174, 252, 249, 184, 78, 133, 75, 233, 253, 100, 239, 16, 33, 121, 194, 113, 84, 69, 70, 141, 82},
			keyLength: 16,
		},
		{
			name:      "AES-192",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey2222212345648",
			cipher:    []byte{136, 74, 208, 187, 3, 44, 189, 63, 1, 173, 123, 112, 202, 111, 4, 144, 25, 71, 78, 42, 139, 168, 127, 194, 83, 104, 196, 93, 49, 172, 254, 137, 147, 255, 153, 22, 171},
			keyLength: 24,
		},
		{
			name:      "AES-256",
			plaintext: "HalloWelt",
			secretKey: "mysecretkey22222mysecretkey22222",
			cipher:    []byte{56, 151, 243, 252, 235, 149, 5, 26, 9, 4, 39, 201, 255, 49, 248, 83, 45, 213, 83, 71, 129, 36, 200, 84, 154, 164, 27, 143, 65, 146, 61, 48, 146, 16, 132, 168, 197},
			keyLength: 32,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			plain, _ := encryption.DecryptGcm([]byte(test.secretKey), test.cipher)
			assert.Equal(t, test.plaintext, string(plain))
		})
	}
}

func TestSha256AES256(t *testing.T) {
	tests := []struct {
		name      string
		plaintext string
		secretKey string
	}{
		{
			plaintext: "HalloWelt",
			secretKey: "Lo4aiph3aikohchei2eeho8aequee3thu7ziegha",
		},
		{
			plaintext: "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
			secretKey: "aoleix1aeJ",
		},
		{
			plaintext: "HalloWelt",
			secretKey: "cah2uleigohbooT2iuqu4luo0HugiK3phu3Meequieh4OhX4aequ7Uy9bei8ohch6Wu4Vah4aewoot4Feighienaekahkay8ayakaegh2ehoo5Eeph3cahm3aiha1yaighohri2giegu",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// encrypt
			cipher, _ := encryption.EncryptSha256AES256([]byte(test.secretKey), []byte(test.plaintext))

			// decrypt
			plain, _ := encryption.DecryptSha256AES256([]byte(test.secretKey), cipher)
			assert.Equal(t, test.plaintext, string(plain))
		})
	}
}
