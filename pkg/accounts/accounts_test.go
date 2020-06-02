package accounts

import (
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestNewAccountant(t *testing.T) {
	// Very basic verify here for now
	accountant := NewAccountant(os.TempDir())
	if accountant == nil {
		t.Error("Failed to create accountant")
	}
}

func TestAccountant_RegisterAccount(t *testing.T) {

	accountant := NewAccountant(os.TempDir())

	// Start by making two accounts

	namea := "one"
	acca, err := accountant.RegisterAccount(namea)
	if err != nil {
		t.Error(err)
	} else if acca.Name != namea {
		t.Errorf("Missmatched account name after register, expected: %s, actual: %s", namea, acca.Name)
	}

	nameb := "two"
	accb, err := accountant.RegisterAccount(nameb)
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
	_, err = accountant.RegisterAccount(namea)
	if err == nil {
		t.Error("Duplicate account name did not produce error")
	}
}

func TestAccountant_LoadSave(t *testing.T) {
	accountant := NewAccountant(os.TempDir())
	if len(accountant.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	name := "one"
	a, err := accountant.RegisterAccount(name)
	if err != nil {
		t.Error(err)
	}

	if len(accountant.Accounts) != 1 {
		t.Error("No new account made")
	} else if accountant.Accounts[a.Id].Name != name {
		t.Error("New account created with wrong name")
	}

	// Save out the accountant
	if err := accountant.Save(); err != nil {
		t.Error(err)
	}

	// Re-create the accountant
	accountant = NewAccountant(os.TempDir())
	if len(accountant.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	// Load the old accountant data
	if err := accountant.Load(); err != nil {
		t.Error(err)
	}

	// Verify we have the same account again
	if len(accountant.Accounts) != 1 {
		t.Error("No account after load")
	} else if accountant.Accounts[a.Id].Name != name {
		t.Error("New account created with wrong name")
	}
}

func TestAccountant_AssignPrimary(t *testing.T) {
	accountant := NewAccountant(os.TempDir())
	if len(accountant.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	name := "one"
	a, err := accountant.RegisterAccount(name)
	if err != nil {
		t.Error(err)
	}

	inst := uuid.New()

	err = accountant.AssignPrimary(a.Id, inst)
	if err != nil {
		t.Error("Failed to set primary for created account")
	} else if accountant.Accounts[a.Id].Primary != inst {
		t.Error("Primary for assigned account is incorrect")
	}
}
