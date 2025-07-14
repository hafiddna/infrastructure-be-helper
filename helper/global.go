package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strings"
)

func JSONPrettyLog(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(string(b))
		log.Println()
		log.Println()
		log.Println()
	}
}

func JSONMarshal(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}

func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type EncryptedData struct {
	IV    string `json:"iv"`
	Value string `json:"value"`
	MAC   string `json:"mac"`
	Tag   string `json:"tag"`
}

func GenerateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize) // AES block size is 16 bytes
	_, err := io.ReadFull(rand.Reader, iv)
	return iv, err
}

func PKCS7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func GenerateHMAC(iv, ciphertext, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(append(iv, ciphertext...))
	return mac.Sum(nil)
}

func EncryptAES256CBC(plaintext, key []byte) (*EncryptedData, error) {
	// Generate IV
	iv, err := GenerateIV()
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad plaintext to match block size
	plaintext = PKCS7Pad(plaintext, aes.BlockSize)

	// Encrypt using CBC mode
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Generate HMAC for integrity
	mac := GenerateHMAC(iv, ciphertext, key)

	// Encode values in base64
	return &EncryptedData{
		IV:    base64.StdEncoding.EncodeToString(iv),
		Value: base64.StdEncoding.EncodeToString(ciphertext),
		MAC:   fmt.Sprintf("%x", mac),
		Tag:   "",
	}, nil
}

func PKCS7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding size")
	}
	return data[:length-padding], nil
}

func VerifyHMAC(iv, ciphertext, key []byte, mac string) bool {
	macGen := hmac.New(sha256.New, key)
	macGen.Write(append(iv, ciphertext...))
	expectedMAC := macGen.Sum(nil)

	// Compare the computed MAC with the provided MAC
	expectedMACHex := fmt.Sprintf("%x", expectedMAC)
	return expectedMACHex == mac
}

func DecryptAES256CBC(encrypted *EncryptedData, key []byte) (string, error) {
	// Decode Base64 values
	iv, err := base64.StdEncoding.DecodeString(encrypted.IV)
	if err != nil {
		return "", fmt.Errorf("error decoding IV: %v", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Value)
	if err != nil {
		return "", fmt.Errorf("error decoding Value: %v", err)
	}

	// Verify HMAC integrity
	if !VerifyHMAC(iv, ciphertext, key, encrypted.MAC) {
		return "", fmt.Errorf("HMAC verification failed, data may be tampered with")
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %v", err)
	}

	// Decrypt using CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = PKCS7Unpad(plaintext)
	if err != nil {
		return "", fmt.Errorf("error removing padding: %v", err)
	}

	return string(plaintext), nil
}

func RoundFloat64(val float64, precision uint) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}

func Like(input, pattern string) bool {
	pattern = regexp.QuoteMeta(pattern)

	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")

	regex := "^" + pattern + "$"

	matched, err := regexp.MatchString(regex, input)
	if err != nil {
		return false
	}
	return matched
}

func ArrayInterfaceContains(arr []interface{}, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func ArrayStringMerge(arrays ...[]string) []string {
	merged := make([]string, 0)
	seen := make(map[string]bool)

	for _, array := range arrays {
		for _, item := range array {
			if !seen[item] {
				seen[item] = true
				merged = append(merged, item)
			}
		}
	}

	return merged
}

func InArrayString(str string, arr []string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}
