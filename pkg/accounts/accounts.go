package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
)

const kDefaultSavePath = "/tmp/accounts.json"

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string `json:"name"`

	// id represents a unique ID per account and is set one registered
	Id uuid.UUID `json:"id"`
}

// Represents the accountant data to store
type accountantData struct {
	Accounts []Account `json:"accounts"`
}

// Accountant manages a set of accounts
type Accountant struct {
	data accountantData
}

// NewAccountant creates a new accountant
func NewAccountant() *Accountant {
	return &Accountant{}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *Accountant) RegisterAccount(acc Account) (Account, error) {

	// Set the account ID to a new UUID
	acc.Id = uuid.New()

	// Verify this acount isn't already registered
	for _, a := range a.data.Accounts {
		if a.Name == acc.Name {
			return Account{}, fmt.Errorf("Account name already registered")
		} else if a.Id == acc.Id {
			return Account{}, fmt.Errorf("Account ID already registered")
		}
	}

	// Simply add the account to the list
	a.data.Accounts = append(a.data.Accounts, acc)

	return acc, nil
}

// Load will load the accountant from data
func (a *Accountant) Load() error {
	// Don't load anything if the file doesn't exist
	_, err := os.Stat(kDefaultSavePath)
	if os.IsNotExist(err) {
		fmt.Printf("File %s didn't exist, loading with fresh accounts data\n", kDefaultSavePath)
		return nil
	}

	if b, err := ioutil.ReadFile(kDefaultSavePath); err != nil {
		return err
	} else if err := json.Unmarshal(b, &a.data); err != nil {
		return err
	}
	return nil
}

// Save will save the accountant data out
func (a *Accountant) Save() error {
	if b, err := json.Marshal(a.data); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(kDefaultSavePath, b, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
