package cryptox_test

import (
	"testing"

	"github.com/tmoeish/gokit/cryptox"
)

func TestMD5(t *testing.T) {
	if cryptox.MD5Str("hello") != "5d41402abc4b2a76b9719d911017c592" {
		t.Fatal("MD5 mismatch")
	}
}

func TestSHA256(t *testing.T) {
	got := cryptox.SHA256Str("hello")
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != want {
		t.Fatalf("SHA256: got %s", got)
	}
}

func TestHMACSHA256(t *testing.T) {
	mac := cryptox.HMACSHA256([]byte("message"), []byte("secret"))
	if len(mac) != 64 {
		t.Fatalf("HMACSHA256 length: %d", len(mac))
	}
	if !cryptox.VerifyHMACSHA256([]byte("message"), []byte("secret"), mac) {
		t.Fatal("VerifyHMACSHA256 failed")
	}
}

func TestBase64(t *testing.T) {
	original := []byte("hello world")
	encoded := cryptox.Base64Encode(original)
	decoded, err := cryptox.Base64Decode(encoded)
	if err != nil || string(decoded) != string(original) {
		t.Fatalf("Base64 roundtrip: %v %v", decoded, err)
	}
}

func TestBase64URL(t *testing.T) {
	original := []byte("hello+world/test")
	encoded := cryptox.Base64URLEncode(original)
	decoded, err := cryptox.Base64URLDecode(encoded)
	if err != nil || string(decoded) != string(original) {
		t.Fatalf("Base64URL roundtrip: %v %v", decoded, err)
	}
}

func TestHex(t *testing.T) {
	data := []byte{0xde, 0xad, 0xbe, 0xef}
	encoded := cryptox.HexEncode(data)
	if encoded != "deadbeef" {
		t.Fatalf("HexEncode: %s", encoded)
	}
	decoded, err := cryptox.HexDecode(encoded)
	if err != nil || string(decoded) != string(data) {
		t.Fatalf("HexDecode: %v %v", decoded, err)
	}
}
