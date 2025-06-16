package logueadores

import (
	"testing"
)

func TestNewTesterLogeador(t *testing.T) {
	claves := []string{"clave1", "clave2"}
	l := NewTesterLogeador(claves)
	if l == nil {
		t.Fatal("Expected non-nil TesterLogeador")
	}
	if len(l.claves) != 2 {
		t.Errorf("Expected 2 claves, got %d", len(l.claves))
	}
}

func TestTesterLogeador_Login_Success(t *testing.T) {
	claves := []string{"user1", "user2"}
	l := NewTesterLogeador(claves)
	if !l.Login("any", "user1") {
		t.Error("Expected Login to return true for valid clave")
	}
	if !l.Login("ignored", "user2") {
		t.Error("Expected Login to return true for valid clave")
	}
}

func TestTesterLogeador_Login_Failure(t *testing.T) {
	claves := []string{"user1", "user2"}
	l := NewTesterLogeador(claves)
	if l.Login("any", "user3") {
		t.Error("Expected Login to return false for invalid clave")
	}
	if l.Login("any", "") {
		t.Error("Expected Login to return false for empty clave")
	}
}

func TestTesterLogeador_Login_EmptyClaves(t *testing.T) {
	l := NewTesterLogeador([]string{})
	if l.Login("any", "user1") {
		t.Error("Expected Login to return false when claves is empty")
	}
}
