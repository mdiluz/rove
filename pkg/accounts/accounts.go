package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
)

const kAccountsFileName = "rove-accounts.json"

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string `json:"name"`

	// Id represents a unique ID per account and is set one registered
	Id uuid.UUID `json:"id"`

	// Primary represents the primary instance that this account owns
	Primary uuid.UUID `json:"primary"`
}

// Represents the accountant data to store
type accountantData struct {
}

// Accountant manages a set of accounts
type Accountant struct {
	Accounts map[uuid.UUID]Account `json:"accounts"`
	dataPath string
}

// NewAccountant creates a new accountant
func NewAccountant(dataPath string) *Accountant {
	return &Accountant{
		dataPath: dataPath,
		Accounts: make(map[uuid.UUID]Account),
	}
}

// RegisterAccount adds an account to the set of internal accounts
func (a *Accountant) RegisterAccount(acc Account) (Account, error) {

	// Set the account ID to a new UUID
	acc.Id = uuid.New()

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

	return acc, nil
}

// path returns the full path to the data file
func (a Accountant) path() string {
	return path.Join(a.dataPath, kAccountsFileName)
}

// Load will load the accountant from data
func (a *Accountant) Load() error {
	// Don't load anything if the file doesn't exist
	_, err := os.Stat(a.path())
	if os.IsNotExist(err) {
		fmt.Printf("File %s didn't exist, loading with fresh accounts data\n", a.path())
		return nil
	}

	if b, err := ioutil.ReadFile(a.path()); err != nil {
		return err
	} else if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	return nil
}

// Save will save the accountant data out
func (a *Accountant) Save() error {
	if b, err := json.MarshalIndent(a, "", "\t"); err != nil {
		return err
	} else {
		if err := ioutil.WriteFile(a.path(), b, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// AssignPrimary assigns primary ownership of an instance to an account
func (a *Accountant) AssignPrimary(account uuid.UUID, instance uuid.UUID) error {

	// Find the account matching the ID
	if this, ok := a.Accounts[account]; ok {
		this.Primary = instance
		a.Accounts[account] = this
	} else {
		return fmt.Errorf("no account found for id: %s", account)
	}

	return nil
}
