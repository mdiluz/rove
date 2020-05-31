package server

import (
	"testing"
)

func TestNewAccountant(t *testing.T) {
	// Very basic verify here for now
	accountant := NewAccountant()
	if accountant == nil {
		t.Error("Failed to create accountant")
	}
}

func TestAccountant_RegisterAccount(t *testing.T) {

	accountant := NewAccountant()

	// Start by making two accounts

	namea := "one"
	a := Account{Name: namea}
	acca, err := accountant.RegisterAccount(a)
	if err != nil {
		t.Error(err)
	} else if acca.Name != namea {
		t.Errorf("Missmatched account name after register, expected: %s, actual: %s", namea, acca.Name)
	}

	nameb := "two"
	b := Account{Name: nameb}
	accb, err := accountant.RegisterAccount(b)
	if err != nil {
		t.Error(err)
	} else if accb.Name != nameb {
		t.Errorf("Missmatched account name after register, expected: %s, actual: %s", nameb, acca.Name)
	}

	// Verify our accounts have differing IDs
	if acca.id == accb.id {
		t.Error("Duplicate account IDs fo separate accounts")
	}

	// Verify another request gets rejected
	_, err = accountant.RegisterAccount(a)
	if err == nil {
		t.Error("Duplicate account name did not produce error")
	}
}
