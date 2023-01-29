package connection_test

import (
	"strings"
	"testing"

	"github.com/ricoschulte/go-myapps/connection"
	"github.com/stretchr/testify/assert"
)

func TestDecryptRc4(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		data     string
		expected string
		err      error
	}{
		{"correct",
			"innovaphoneAppClient:usr:5bba709c4b68374d:ip411",
			"568665fc2f134d3133325b1f30d66e4e6c9f6f9d165e2f5ce7bbfff1b1de68da",
			"9ae4b9193f1362019e19009033400109",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := connection.DecryptRc4(test.key, test.data)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestHex2bin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"A", "568665fc2f134d3133325b1f30d66e4e6c9f6f9d165e2f5ce7bbfff1b1de68da", []int{
			86,
			134,
			101,
			252,
			47,
			19,
			77,
			49,
			51,
			50,
			91,
			31,
			48,
			214,
			110,
			78,
			108,
			159,
			111,
			157,
			22,
			94,
			47,
			92,
			231,
			187,
			255,
			241,
			177,
			222,
			104,
			218,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := connection.Hex2bin(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetRandomHexString(t *testing.T) {
	testCases := []struct {
		n        int
		expected int
	}{
		{5, 10},
		{6, 12},
		{7, 14},
		{8, 16},
		{9, 18},
		{10, 20},
	}
	generated := make(map[string]bool)

	for _, tc := range testCases {
		for i := 0; i < 100; i++ {

			str := connection.GetRandomHexString(tc.n)
			if len(str) != tc.expected {
				t.Errorf("Expected length %d, but got %d", tc.expected, len(str))
			}
			validChars := "abcdef0123456789"
			for _, c := range str {
				if !strings.ContainsRune(validChars, c) {
					t.Errorf("Invalid character %c found in string %s", c, str)
				}
			}

			if generated[str] {
				t.Errorf("Duplicate string %s generated, could happen some times. #TODO find a better way to test randomness", str)
			}
			generated[str] = true
		}
	}
}
