package mongohelpers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	controller "github.com/skyakashh/revpay/controllers"
	"github.com/skyakashh/revpay/models"
	"go.mongodb.org/mongo-driver/bson"
)

//for creating user

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	if !username_verify(user) {
		log.Fatal("not a unique username")
	}
	createdUser := createUser(user)
	json.NewEncoder(w).Encode(createdUser)

}

func username_verify(user models.User) bool {
	var userVerify models.User
	err := controller.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&userVerify)
	return err == nil
}

func createUser(user models.User) models.User {
	_, err := controller.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return user
}

// for authenticating password and username

func authentication(user models.UserAuth) bool {
	var userVerify models.User
	err := controller.Collection.FindOne(context.TODO(), bson.M{"username": user.Username, "password": user.Password}).Decode(&userVerify)
	return err != nil
}

// for creating account

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Response
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	created := account(user)

	json.NewEncoder(w).Encode(created)
}

func account(user models.Response) models.Account {

	if len(user.BankAccount) > 10 || len(user.IFSC) > 8 {
		log.Fatal("incorrect credentials")
	}
	var userVerify models.User
	err := controller.Collection.FindOne(context.TODO(), bson.M{"username": user.Username, "password": user.Password}).Decode(&userVerify)
	if err != nil {
		log.Fatal("user not found")
	}
	var account models.Account
	account.AccountID = user.AccountID
	account.BankAccount = user.BankAccount
	account.DailyLimit = 1000
	account.CurrentBalance = 0
	account.IFSC = user.IFSC
	account.Status = user.Status
	account.AllowCredit = user.AllowCredit
	account.AllowDebit = user.AllowDebit

	userVerify.Accounts = append(userVerify.Accounts, account)

	err1 := controller.Collection.FindOneAndReplace(context.TODO(), bson.M{"username": user.Username, "password": user.Password}, userVerify)
	if err1 != nil {
		log.Fatal(err1)
	}
	return account

}

// for withdrawl money
