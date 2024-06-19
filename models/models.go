package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	BankAccount    string             `json:"bankaccount" validate:"required,min=2,max=8"`
	IFSC           string             `json:"ifsc" validate:"required,min=2,max=8"`
	Status         string             `json:"status,omitempty"`
	AllowCredit    bool               `json:"allowcredit,omitempty"`
	AllowDebit     bool               `json:"allowdebit,omitempty"`
	DailyLimit     float64            `json:"dailylimit,omitempty"`
	CurrentBalance float64            `json:"currentbalance,omitempty"`
	LastUpdated    time.Time          `json:"time,omitempty"`
}

type UserAuth struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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

type Balance struct {
	BankAccount string `json:"bankaccount"`
	IFSC        string `json:"ifsc"`
}
