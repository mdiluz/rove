package accounts

import (
	"testing"

	"github.com/google/uuid"
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

	namea := uuid.New().String()
	acca, err := accountant.RegisterAccount(namea)
	if err != nil {
		t.Error(err)
	} else if acca.Name != namea {
		t.Errorf("Missmatched account name after register, expected: %s, actual: %s", namea, acca.Name)
	}

	nameb := uuid.New().String()
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

func TestAccountant_AssignGetRover(t *testing.T) {
	accountant := NewAccountant()
	if len(accountant.Accounts) != 0 {
		t.Error("New accountant created with non-zero account number")
	}

	name := uuid.New().String()
	a, err := accountant.RegisterAccount(name)
	if err != nil {
		t.Error(err)
	}

	inst := uuid.New()

	err = accountant.AssignRover(a.Id, inst)
	if err != nil {
		t.Error("Failed to set rover for created account")
	} else if accountant.Accounts[a.Id].Rover != inst {
		t.Error("Rover for assigned account is incorrect")
	} else if id, err := accountant.GetRover(a.Id); err != nil {
		t.Error("Failed to get rover for account")
	} else if id != inst {
		t.Error("Fetched rover is incorrect for account")
	}
}
