// Package cryptox provides cryptographic utility functions.
// All functions in this package use only stdlib crypto; no external deps.
package cryptox

import (
	"crypto/hmac"
	"crypto/md5"  //nolint:gosec // MD5 used for non-security checksums
	"crypto/sha1" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"hash"
)

// MD5 returns the MD5 hex digest of data.
// NOTE: MD5 is cryptographically broken; use for checksums only, not security.
func MD5(data []byte) string {
	h := md5.Sum(data) //nolint:gosec
	return hex.EncodeToString(h[:])
}

// MD5Str returns the MD5 hex digest of s.
func MD5Str(s string) string {
	return MD5([]byte(s))
}

// SHA1 returns the SHA-1 hex digest of data.
// NOTE: SHA-1 is deprecated for security; use SHA-256 or SHA-512 instead.
func SHA1(data []byte) string {
	h := sha1.Sum(data) //nolint:gosec
	return hex.EncodeToString(h[:])
}

// SHA256 returns the SHA-256 hex digest of data.
func SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// SHA256Str returns the SHA-256 hex digest of s.
func SHA256Str(s string) string {
	return SHA256([]byte(s))
}

// SHA512 returns the SHA-512 hex digest of data.
func SHA512(data []byte) string {
	h := sha512.Sum512(data)
	return hex.EncodeToString(h[:])
}

// SHA512Str returns the SHA-512 hex digest of s.
func SHA512Str(s string) string {
	return SHA512([]byte(s))
}

// HMACSHA256 returns the HMAC-SHA256 hex digest of data with the given key.
func HMACSHA256(data, key []byte) string {
	return hexHMAC(sha256.New, data, key)
}

// HMACSHA512 returns the HMAC-SHA512 hex digest of data with the given key.
func HMACSHA512(data, key []byte) string {
	return hexHMAC(sha512.New, data, key)
}

// VerifyHMACSHA256 reports whether mac matches the HMAC-SHA256 of data with key.
func VerifyHMACSHA256(data, key []byte, mac string) bool {
	expected := HMACSHA256(data, key)
	return hmac.Equal([]byte(expected), []byte(mac))
}

// Base64Encode encodes data to standard base64.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode decodes a standard base64 string.
func Base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Base64URLEncode encodes data to URL-safe base64 (no padding).
func Base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Base64URLDecode decodes a URL-safe base64 string (no padding).
func Base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// HexEncode returns the hex encoding of data.
func HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode decodes a hex string.
func HexDecode(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func hexHMAC(newHash func() hash.Hash, data, key []byte) string {
	mac := hmac.New(newHash, key)
	mac.Write(data) //nolint:errcheck
	return hex.EncodeToString(mac.Sum(nil))
}
