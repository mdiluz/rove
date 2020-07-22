package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SimpleAccountant manages a set of accounts
type SimpleAccountant struct {
	Accounts map[string]Account
}

// NewSimpleAccountant creates a new accountant
func NewSimpleAccountant() Accountant {
	return &SimpleAccountant{
		Accounts: make(map[string]Account),
	}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *SimpleAccountant) RegisterAccount(name string) (acc Account, err error) {

	// Set up the account info
	acc.Name = name
	acc.Data = make(map[string]string)

	// Verify this acount isn't already registered
	for _, a := range a.Accounts {
		if a.Name == acc.Name {
			return Account{}, fmt.Errorf("account name already registered: %s", a.Name)
		}
	}

	// Set the creation time
	acc.Data["created"] = time.Now().String()

	// Create a secret
	acc.Data["secret"] = uuid.New().String()

	// Simply add the account to the map
	a.Accounts[acc.Name] = acc

	return
}

// VerifySecret verifies if an account secret is correct
func (a *SimpleAccountant) VerifySecret(account string, secret string) (bool, error) {
	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		return this.Data["secret"] == secret, nil
	}

	return false, fmt.Errorf("no account found for id: %s", account)
}

// GetSecret gets the internal secret
func (a *SimpleAccountant) GetSecret(account string) (string, error) {
	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		return this.Data["secret"], nil
	}

	return "", fmt.Errorf("no account found for id: %s", account)
}

// AssignData assigns data to an account
func (a *SimpleAccountant) AssignData(account string, key string, value string) error {

	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		this.Data[key] = value
		a.Accounts[account] = this
	} else {
		return fmt.Errorf("no account found for id: %s", account)
	}

	return nil
}

// GetValue gets the rover rover for the account
func (a *SimpleAccountant) GetValue(account string, key string) (string, error) {
	// Find the account matching the ID
	this, ok := a.Accounts[account]
	if !ok {
		return "", fmt.Errorf("no account found for id: %s", account)
	}
	return this.Data[key], nil
}
