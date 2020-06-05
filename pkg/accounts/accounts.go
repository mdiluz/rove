package accounts

import (
	"fmt"

	"github.com/google/uuid"
)

const kAccountsFileName = "rove-accounts.json"

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string `json:"name"`

	// Id represents a unique ID per account and is set one registered
	Id uuid.UUID `json:"id"`

	// Rover represents the rover that this account owns
	Rover uuid.UUID `json:"rover"`
}

// Represents the accountant data to store
type accountantData struct {
}

// Accountant manages a set of accounts
type Accountant struct {
	Accounts map[uuid.UUID]Account `json:"accounts"`
}

// NewAccountant creates a new accountant
func NewAccountant() *Accountant {
	return &Accountant{
		Accounts: make(map[uuid.UUID]Account),
	}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *Accountant) RegisterAccount(name string) (acc Account, err error) {

	// Set the account ID to a new UUID
	acc.Id = uuid.New()
	acc.Name = name

	// Verify this acount isn't already registered
	for _, a := range a.Accounts {
		if a.Name == acc.Name {
			return Account{}, fmt.Errorf("Account name already registered")
		} else if a.Id == acc.Id {
			return Account{}, fmt.Errorf("Account ID already registered")
		}
	}

	// Simply add the account to the map
	a.Accounts[acc.Id] = acc

	return
}

// AssignRover assigns rover ownership of an rover to an account
func (a *Accountant) AssignRover(account uuid.UUID, rover uuid.UUID) error {

	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		this.Rover = rover
		a.Accounts[account] = this
	} else {
		return fmt.Errorf("no account found for id: %s", account)
	}

	return nil
}

// GetRover gets the rover rover for the account
func (a *Accountant) GetRover(account uuid.UUID) (uuid.UUID, error) {
	// Find the account matching the ID
	if this, ok := a.Accounts[account]; !ok {
		return uuid.UUID{}, fmt.Errorf("no account found for id: %s", account)
	} else if this.Rover == uuid.Nil {
		return uuid.UUID{}, fmt.Errorf("no rover spawned for account %s", account)
	} else {
		return this.Rover, nil
	}
}
