package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// func PrettyPrintJSON(data []map[string]string) string {
func PrettyPrintJSON(data any) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
	return string(b)
}

/*
splits a string by delimiter but ignores delimiter escaped by escape

in:  aaaaa:bbbb:cccc\:dddd
out: [aaaaa,bbbb,cccc\:dddd]

splitIgnoreEscape(line2, ':', '\\')
*/
func SplitIgnoreEscape(str string, delimiter byte, escape byte) []string {
	var parts []string
	var buffer strings.Builder

	escaped := false

	for i := 0; i < len(str); i++ {
		char := str[i]

		if !escaped && char == escape {
			escaped = true
			continue
		}

		if escaped {
			if char != delimiter && char != escape {
				buffer.WriteByte(escape)
			}
			buffer.WriteByte(char)
			escaped = false
			continue
		}

		if char == delimiter {
			parts = append(parts, buffer.String())
			buffer.Reset()
			continue
		}

		buffer.WriteByte(char)
	}

	parts = append(parts, buffer.String())

	return parts
}

/*
splits a string by delimiter and returns a slice of non empty strings

SplitStringIgnoreEmpty(value, ";")
in "aaaa;bbbb;;cccc"
out ["aaaa", "bbbb", "cccc"]
*/
func SplitStringIgnoreEmpty(value string, delimiter string) []string {
	filteredSlice := make([]string, 0)
	for _, str := range strings.Split(value, ";") {
		if strings.TrimSpace(str) != "" {
			filteredSlice = append(filteredSlice, str)
		}
	}
	return filteredSlice
}

// sha256 hash from map[string]string
func Sha265HashFromMap(m map[string]string) string {
	// Get a sorted slice of the keys
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Concatenate the key-value pairs
	var str string
	for _, key := range keys {
		str += key + m[key]
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(str))

	// Return the hexadecimal encoding of the hash
	return hex.EncodeToString(hash[:])
}

/*
takes two []string slices as input and returns true if all the elements in the first slice (s1) are present in the second slice (s2), regardless of their order, and false otherwise.

The function first checks if the two slices have the same length, and if not, it immediately returns false. Then, it creates a map (m) and adds each element of the second slice (s2) to the map as a key with a value of true.

Next, it iterates over the elements of the first slice (s1) and checks if each element is present in the map (m). If any element is not present in the map, the function returns false. If all elements are present in the map, the function returns true.
*/
func AreEqualSlices(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	m := make(map[string]bool)
	for _, v := range s2 {
		m[v] = true
	}

	for _, v := range s1 {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}
