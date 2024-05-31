package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username string               `json:"username"`
	Password string               `json:"password"`
	Accounts []primitive.ObjectID `json:"accounts,omitempty"`
}

type Account struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username"`
	BankAccount    string             `json:"bankaccount"`
	IFSC           string             `json:"ifsc"`
	Status         string             `json:"status,omitempty"`
	AllowCredit    bool               `json:"allow_credit,omitempty"`
	AllowDebit     bool               `json:"allow_debit,omitempty"`
	DailyLimit     float64            `json:"daily_limit,omitempty"`
	CurrentBalance float64            `json:"current_balance,omitempty"`
	LastUpdated    time.Time          `json:"time,omitempty"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Withdraw struct {
	BankAccount string  `json:"bankaccount"`
	IFSC        string  `json:"ifsc"`
	Amount      float64 `json:"amount"`
}

type UserId struct {
	Username  string `json:"username"`
	AccountId string `json:"account_id"`
}

type Deposit struct {
	BankAccount string  `json:"bankaccount"`
	IFSC        string  `json:"ifsc"`
	PaymentMode string  `json:"paymentmode"`
	Amount      float64 `json:"amount"`
}
