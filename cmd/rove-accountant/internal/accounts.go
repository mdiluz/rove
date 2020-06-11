package internal

import (
	"fmt"
	"time"
)

const kAccountsFileName = "rove-accounts.json"

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string `json:"name"`

	// Data represents internal account data
	Data map[string]string `json:"data"`
}

// Represents the accountant data to store
type accountantData struct {
}

// Accountant manages a set of accounts
type Accountant struct {
	Accounts map[string]Account `json:"accounts"`
}

// NewAccountant creates a new accountant
func NewAccountant() *Accountant {
	return &Accountant{
		Accounts: make(map[string]Account),
	}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *Accountant) RegisterAccount(name string) (acc Account, err error) {

	// Set the account name
	acc.Name = name
	acc.Data = make(map[string]string)

	// Verify this acount isn't already registered
	for _, a := range a.Accounts {
		if a.Name == acc.Name {
			return Account{}, fmt.Errorf("account name already registered")
		}
	}

	// Set the creation time
	acc.Data["created"] = time.Now().String()

	// Simply add the account to the map
	a.Accounts[acc.Name] = acc

	return
}

// AssignRover assigns data to an account
func (a *Accountant) AssignData(account string, key string, value string) error {

	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		this.Data[key] = value
		a.Accounts[account] = this
	} else {
		return fmt.Errorf("no account found for id: %s", account)
	}

	return nil
}

// GetRover gets the rover rover for the account
func (a *Accountant) GetValue(account string, key string) (string, error) {
	// Find the account matching the ID
	if this, ok := a.Accounts[account]; !ok {
		return "", fmt.Errorf("no account found for id: %s", account)
	} else {
		return this.Data[key], nil
	}
}
