package models

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	AccountID      string  `json:"account_id"`
	BankAccount    string  `json:"bank_account"`
	IFSC           string  `json:"ifsc"`
	Status         string  `json:"status"`
	AllowCredit    bool    `json:"allow_credit"`
	AllowDebit     bool    `json:"allow_debit"`
	DailyLimit     float64 `json:"daily_limit"`
	CurrentBalance float64 `json:"current_balance"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccountID   string `json:"account_id"`
	BankAccount string `json:"bank_account"`
	IFSC        string `json:"ifsc"`
	Status      string `json:"status"`
	AllowCredit bool   `json:"allow_credit"`
	AllowDebit  bool   `json:"allow_debit"`
}

type Withdraw struct {
	Username    string  `json:"username"`
	BankAccount string  `json:"bank_account"`
	IFSC        string  `json:"ifsc"`
	Amount      float64 `json:"amount"`
}

type UserId struct {
	Username  string `json:"username"`
	AccountId string `json:"account_id"`
}

type Deposit struct {
	Username    string  `json:"username"`
	BankAccount string  `json:"bank_account"`
	IFSC        string  `json:"ifsc"`
	PaymentMode string  `json:"payment_mode"`
	Amount      float64 `json:"amount"`
}
