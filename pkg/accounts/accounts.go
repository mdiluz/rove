package accounts

// Accountant decribes something that stores accounts and account values
type Accountant interface {
	// RegisterAccount will register a new account and return it's info
	RegisterAccount(name string) (acc Account, err error)

	// AssignData stores a custom account key value pair
	AssignData(account string, key string, value string) error

	// GetValue returns custom account data for a specific key
	GetValue(account string, key string) (string, error)

	// VerifySecret will verify whether the account secret matches with the
	VerifySecret(account string, secret string) (bool, error)

	// GetSecret gets the secret associated with an account
	GetSecret(account string) (string, error)
}

// Account represents a registered user
type Account struct {
	// Name simply describes the account and must be unique
	Name string `json:"name"`

	// Data represents internal account data
	Data map[string]string `json:"data"`
}
