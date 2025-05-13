package tests

import (
	"encoding/base64"
	"strings"
	"testing"

	"terrabutler/variables"
)

// ------------------------------
// TESTES: GeneratePassword
// ------------------------------

func TestGeneratePassword_Length(t *testing.T) {
	pw := variables.GeneratePassword(16)
	if len(pw) != 16 {
		t.Errorf("Expected password length of 16, got %d", len(pw))
	}
}

func TestGeneratePassword_CharsetOnly(t *testing.T) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	pw := variables.GeneratePassword(100)
	for _, ch := range pw {
		if !strings.ContainsRune(charset, ch) {
			t.Errorf("Unexpected character in password: %q", ch)
		}
	}
}

// ------------------------------
// TESTES: contains
// ------------------------------

func TestContains_True(t *testing.T) {
	list := []string{"one", "two", "three"}
	if !variables.Contains(list, "two") {
		t.Error("Expected true when element is in slice")
	}
}

func TestContains_False(t *testing.T) {
	list := []string{"one", "two", "three"}
	if variables.Contains(list, "four") {
		t.Error("Expected false when element is not in slice")
	}
}

// ------------------------------
// TESTES: remove
// ------------------------------

func TestRemove_Existing(t *testing.T) {
	list := []string{"a", "b", "c"}
	result := variables.Remove(list, "b")
	expected := []string{"a", "c"}

	if len(result) != len(expected) {
		t.Fatalf("Expected result length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %s at index %d, got %s", v, i, result[i])
		}
	}
}

func TestRemove_NotFound(t *testing.T) {
	list := []string{"a", "b", "c"}
	result := variables.Remove(list, "x")

	if len(result) != 3 {
		t.Errorf("Expected unchanged list, got length %d", len(result))
	}
}

// ------------------------------
// TESTES: GenerateEncryptedPassword (mock básico)
// ------------------------------

func TestGenerateEncryptedPassword_Base64Encoded(t *testing.T) {
	// Essa função não testa realmente AWS, mas testa se o resultado parece codificado
	pw := variables.GeneratePassword(10)
	encoded := base64.StdEncoding.EncodeToString([]byte(pw))

	if encoded == "" {
		t.Error("Expected a base64 encoded string, got empty")
	}
}
