package accounts

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
	if acca.Id == accb.Id {
		t.Error("Duplicate account IDs fo separate accounts")
	}

	// Verify another request gets rejected
	_, err = accountant.RegisterAccount(a)
	if err == nil {
		t.Error("Duplicate account name did not produce error")
	}
}

func TestAccountant_LoadSave(t *testing.T) {
	accountant := NewAccountant()
	if len(accountant.data.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	name := "one"
	a := Account{Name: name}
	a, err := accountant.RegisterAccount(a)
	if err != nil {
		t.Error(err)
	}

	if len(accountant.data.Accounts) != 1 {
		t.Error("No new account made")
	} else if accountant.data.Accounts[0].Name != name {
		t.Error("New account created with wrong name")
	}

	// Save out the accountant
	if err := accountant.Save(); err != nil {
		t.Error(err)
	}

	// Re-create the accountant
	accountant = NewAccountant()
	if len(accountant.data.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	// Load the old accountant data
	if err := accountant.Load(); err != nil {
		t.Error(err)
	}

	// Verify we have the same account again
	if len(accountant.data.Accounts) != 1 {
		t.Error("No account after load")
	} else if accountant.data.Accounts[0].Name != name {
		t.Error("New account created with wrong name")
	}
}
