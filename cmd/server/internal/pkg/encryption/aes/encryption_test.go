package aes

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptAES_OK(t *testing.T) {
	t.Parallel()

	keys := [][]byte{
		bytes.Repeat([]byte{1}, 16),
		bytes.Repeat([]byte{2}, 24),
		bytes.Repeat([]byte{3}, 32),
	}

	plaintexts := [][]byte{
		[]byte("hello"),
		[]byte("hello world"),
		bytes.Repeat([]byte("a"), 16),
		bytes.Repeat([]byte("b"), 31),
		bytes.Repeat([]byte("c"), 100),
	}

	for _, key := range keys {
		for _, pt := range plaintexts {

			ciphertext, err := EncryptAES(pt, key)
			if err != nil {
				t.Fatalf("EncryptAES failed: %v", err)
			}

			if bytes.Equal(ciphertext, pt) {
				t.Fatalf("ciphertext equals plaintext (encryption ineffective)")
			}

			decrypted, err := DecryptAES(ciphertext, key)
			if err != nil {
				t.Fatalf("DecryptAES failed: %v", err)
			}

			if !bytes.Equal(decrypted, pt) {
				t.Fatalf("decrypted != plaintext\nwant=%q\ngot =%q", pt, decrypted)
			}
		}
	}
}

func TestEncryptAES_InvalidKeyLength(t *testing.T) {
	t.Parallel()

	invalidKeys := [][]byte{
		{},
		[]byte("short"),
		bytes.Repeat([]byte{1}, 15),
		bytes.Repeat([]byte{1}, 17),
		bytes.Repeat([]byte{1}, 31),
	}

	for _, key := range invalidKeys {
		_, err := EncryptAES([]byte("data"), key)
		if err == nil {
			t.Fatalf("expected error for key len=%d", len(key))
		}
	}
}

func TestDecryptAES_InvalidKeyLength(t *testing.T) {
	t.Parallel()

	key := []byte("short")
	_, err := DecryptAES([]byte("ciphertext"), key)
	if err == nil {
		t.Fatalf("expected error for invalid key length")
	}
}

func TestDecryptAES_TooShortCiphertext(t *testing.T) {
	t.Parallel()

	key := bytes.Repeat([]byte{1}, 16)

	_, err := DecryptAES([]byte("short"), key)
	if err == nil {
		t.Fatalf("expected error for short ciphertext")
	}
}

func TestEncryptAES_RandomIV(t *testing.T) {
	t.Parallel()

	key := bytes.Repeat([]byte{1}, 16)
	plaintext := []byte("same plaintext")

	c1, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES failed: %v", err)
	}

	c2, err := EncryptAES(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptAES failed: %v", err)
	}

	if bytes.Equal(c1, c2) {
		t.Fatalf("ciphertexts are equal, IV is not random")
	}
}

func TestDecryptAES_WrongKey(t *testing.T) {
	t.Parallel()

	key1 := bytes.Repeat([]byte{1}, 16)
	key2 := bytes.Repeat([]byte{2}, 16)

	plaintext := []byte("secret data")

	ciphertext, err := EncryptAES(plaintext, key1)
	if err != nil {
		t.Fatalf("EncryptAES failed: %v", err)
	}

	decrypted, err := DecryptAES(ciphertext, key2)
	if err == nil && bytes.Equal(decrypted, plaintext) {
		t.Fatalf("decryption with wrong key should not succeed")
	}
}
