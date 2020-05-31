package server

import (
	"fmt"

	"github.com/google/uuid"
)

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string

	// id represents a unique ID per account and is set one registered
	id uuid.UUID
}

// Accountant manages a set of accounts
type Accountant struct {
	accounts []Account
}

// NewAccountant creates a new accountant
func NewAccountant() *Accountant {
	return &Accountant{}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *Accountant) RegisterAccount(acc Account) (Account, error) {

	// Set the account ID to a new UUID
	acc.id = uuid.New()

	// Verify this acount isn't already registered
	for _, a := range a.accounts {
		if a.Name == acc.Name {
			return Account{}, fmt.Errorf("Account name already registered")
		} else if a.id == acc.id {
			return Account{}, fmt.Errorf("Account ID already registered")
		}
	}

	// Simply add the account to the list
	a.accounts = append(a.accounts, acc)

	return acc, nil
}
