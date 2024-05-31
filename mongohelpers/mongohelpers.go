package mongohelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	controller "github.com/skyakashh/revpay/controllers"
	"github.com/skyakashh/revpay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	fmt.Println(userVerify)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document was found
			fmt.Println("No document found")
			return true
		}
		// Handle other potential errors
		log.Fatal(err)
		return false
	}
	return false
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

func account(user models.Response) models.User {

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

	filter := bson.M{"username": user.Username, "password": user.Password}

	err = controller.Collection.FindOneAndReplace(context.TODO(), filter, userVerify).Decode(&userVerify)
	if err != nil {
		log.Fatal(err)
	}
	return userVerify

}

// TODO: for withdrawl money

func Withdrawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Withdraw
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	// finding user in database
	var userVerify models.User
	err = controller.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&userVerify)
	if err != nil {
		log.Fatal("user not found")
	}
	fmt.Println(userVerify)
	for i, value := range userVerify.Accounts {
		if value.BankAccount == user.BankAccount && user.IFSC == value.IFSC {
			if value.DailyLimit < user.Amount {
				log.Fatal("limit exceede")
			}
			// make changes in balance
			userVerify.Accounts[i].CurrentBalance -= user.Amount
			//daily limit
			userVerify.Accounts[i].DailyLimit -= user.Amount

			err = controller.Collection.FindOneAndReplace(context.TODO(), bson.M{"username": user.Username}, userVerify).Decode(&userVerify)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(userVerify)
			return
		}
	}

}

// for depositing the money

func Deposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Deposit
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	// finding user in database
	var userVerify models.User
	err = controller.Collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&userVerify)
	if err != nil {
		log.Fatal("user not found")
	}
	fmt.Println(userVerify)
	for i, value := range userVerify.Accounts {
		if value.BankAccount == user.BankAccount && user.IFSC == value.IFSC {
			if value.DailyLimit < user.Amount {
				log.Fatal("limit exceeded")
			}
			if user.PaymentMode == "debit" && !value.AllowDebit {
				log.Fatal("debit payment mode not allowed for this account")
			}
			if user.PaymentMode == "credit" && !value.AllowCredit {
				log.Fatal("credit payment mode not allowed for this account")
			}
			if user.PaymentMode != "credit" && user.PaymentMode != "debit" {
				log.Fatal("invalid payment mode")
			}
			// make changes in balance
			userVerify.Accounts[i].CurrentBalance += user.Amount
			//daily limit
			userVerify.Accounts[i].DailyLimit -= user.Amount

			err = controller.Collection.FindOneAndReplace(context.TODO(), bson.M{"username": user.Username}, userVerify).Decode(&userVerify)
			if err != nil {
				log.Fatal(err)
			}
			json.NewEncoder(w).Encode(userVerify)
			return
		}
	}

}
