package user

import (
	"strings"
	"testing"
)

func TestHashString_FormatAndVerify(t *testing.T) {
	t.Parallel()

	pass := "P@ssw0rd-123"

	encoded, err := HashString(pass)
	if err != nil {
		t.Fatalf("HashString error: %v", err)
	}

	if !strings.HasPrefix(encoded, "$argon2id$") {
		t.Fatalf("expected prefix $argon2id$, got: %q", encoded)
	}
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		t.Fatalf("expected 6 parts split by '$', got %d: %v", len(parts), parts)
	}
	if parts[1] != "argon2id" {
		t.Fatalf("expected parts[1]=argon2id, got %q", parts[1])
	}
	if !strings.HasPrefix(parts[2], "v=") {
		t.Fatalf("expected version part 'v=..', got %q", parts[2])
	}
	if !strings.HasPrefix(parts[3], "m=") || !strings.Contains(parts[3], "t=") || !strings.Contains(parts[3], "p=") {
		t.Fatalf("expected params part like 'm=..,t=..,p=..', got %q", parts[3])
	}
	if parts[4] == "" || parts[5] == "" {
		t.Fatalf("expected non-empty salt and hash parts")
	}

	ok, err := VerifyString(pass, encoded)
	if err != nil {
		t.Fatalf("VerifyString error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true for correct password")
	}

	ok, err = VerifyString("wrong-password", encoded)
	if err != nil {
		t.Fatalf("VerifyString error: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for wrong password")
	}
}

func TestHashString_SamePasswordDifferentHashes(t *testing.T) {
	t.Parallel()

	pass := "same-pass"

	h1, err := HashString(pass)
	if err != nil {
		t.Fatalf("HashString error: %v", err)
	}
	h2, err := HashString(pass)
	if err != nil {
		t.Fatalf("HashString error: %v", err)
	}

	if h1 == h2 {
		t.Fatalf("expected different hashes for same password (salted), but got equal")
	}

	ok, err := VerifyString(pass, h1)
	if err != nil || !ok {
		t.Fatalf("expected verify ok for h1, ok=%v err=%v", ok, err)
	}
	ok, err = VerifyString(pass, h2)
	if err != nil || !ok {
		t.Fatalf("expected verify ok for h2, ok=%v err=%v", ok, err)
	}
}

func TestDecodeHash_OK(t *testing.T) {
	t.Parallel()

	pass := "decode-me"
	encoded, err := HashString(pass)
	if err != nil {
		t.Fatalf("HashString error: %v", err)
	}

	p, salt, hash, err := decodeHash(encoded)
	if err != nil {
		t.Fatalf("decodeHash error: %v", err)
	}
	if p == nil || len(salt) == 0 || len(hash) == 0 {
		t.Fatalf("expected non-empty parsed fields")
	}
	if p.memory == 0 || p.iterations == 0 || p.parallelism == 0 {
		t.Fatalf("expected params to be set, got: %+v", p)
	}
	if p.saltLength != uint32(len(salt)) {
		t.Fatalf("expected saltLength=%d, got %d", len(salt), p.saltLength)
	}
	if p.keyLength != uint32(len(hash)) {
		t.Fatalf("expected keyLength=%d, got %d", len(hash), p.keyLength)
	}
}

func TestDecodeHash_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, _, err := decodeHash("not-a-valid-hash")
	if err == nil {
		t.Fatalf("expected error for invalid format")
	}
	if !errorsIs(err, ErrInvalidHash) {
		t.Fatalf("expected ErrInvalidHash, got: %v", err)
	}
}

func TestDecodeHash_IncompatibleVersion(t *testing.T) {
	t.Parallel()

	bad := "$argon2id$v=18$m=65536,t=3,p=1$YWJjZGVmZ2hpamtsbW5vcA$YWJjZGVmZ2hpamtsbW5vcA"

	_, _, _, err := decodeHash(bad)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errorsIs(err, ErrIncompatibleVersion) {
		t.Fatalf("expected ErrIncompatibleVersion, got: %v", err)
	}
}

func TestDecodeHash_BadParams(t *testing.T) {
	t.Parallel()

	bad := "$argon2id$v=19$m=NOPE,t=3,p=1$YWJjZGVmZ2hpamtsbW5vcA$YWJjZGVmZ2hpamtsbW5vcA"

	_, _, _, err := decodeHash(bad)
	if err == nil {
		t.Fatalf("expected error for bad params")
	}
}

func TestDecodeHash_BadSaltBase64(t *testing.T) {
	t.Parallel()

	bad := "$argon2id$v=19$m=65536,t=3,p=1$%%%NOTBASE64%%%$YWJjZGVmZ2hpamtsbW5vcA"
	_, _, _, err := decodeHash(bad)
	if err == nil {
		t.Fatalf("expected error for bad salt base64")
	}
}

func TestDecodeHash_BadHashBase64(t *testing.T) {
	t.Parallel()

	bad := "$argon2id$v=19$m=65536,t=3,p=1$YWJjZGVmZ2hpamtsbW5vcA$%%%NOTBASE64%%%"
	_, _, _, err := decodeHash(bad)
	if err == nil {
		t.Fatalf("expected error for bad hash base64")
	}
}

func TestGenerateRandom_LengthAndNotAllSame(t *testing.T) {
	t.Parallel()

	b1, err := generateRandom(16)
	if err != nil {
		t.Fatalf("generateRandom error: %v", err)
	}
	if len(b1) != 16 {
		t.Fatalf("expected len=16, got %d", len(b1))
	}

	b2, err := generateRandom(16)
	if err != nil {
		t.Fatalf("generateRandom error: %v", err)
	}
	if len(b2) != 16 {
		t.Fatalf("expected len=16, got %d", len(b2))
	}

	if strings.EqualFold(string(b1), string(b2)) && string(b1) == string(b2) {
		t.Fatalf("unexpected: two random outputs are equal")
	}
}

func errorsIs(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	return err == target
}
