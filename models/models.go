package models

type User struct {
	Username string    `json:"username"`
	Password string    `json:"password"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	UserID         uint    `json:"user_id"`
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
	UserID      uint   `json:"user_id"`
	AccountID   string `json:"account_id"`
	BankAccount string `json:"bank_account"`
	IFSC        string `json:"ifsc"`
	Status      string `json:"status"`
	AllowCredit bool   `json:"allow_credit"`
	AllowDebit  bool   `json:"allow_debit"`
}
