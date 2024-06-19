package mongohelpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	controller "github.com/skyakashh/revpay/controllers"
	"github.com/skyakashh/revpay/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var hmacSampleSecret = []byte("lqdbgouegwucb")

//for creating user

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request recieved")
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	if !username_verify(user) {
		log.Fatal("not a unique username")
	}

	// hashing the password
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// user.Password = string(hashedPassword)
	// if err != nil {
	// 	http.Error(w, "Error hashing password", http.StatusInternalServerError)
	// 	return
	// }

	// creating the user
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

func Authentication(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(user)
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// user.Password = string(hashedPassword)

	//fmt.Println(string(hashedPassword))
	var authuser models.User
	filter := bson.M{"username": user.Username, "password": user.Password}
	err = controller.Collection.FindOne(context.TODO(), filter).Decode(&authuser)
	// return err != nil
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(authuser)
	json.NewEncoder(w).Encode(authuser)
	// verifying the password

	// expirationTime := time.Now().Add(5 * time.Minute)
	// claims := &models.UserAuth{
	// 	Username: user.Username,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(expirationTime),
	// 	},
	// }

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// fmt.Println(token)
	// tokenString, err := token.SignedString(hmacSampleSecret)
	// fmt.Println(tokenString, err)

	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })
}

// Sign and get the complete encoded token as a string using the secret

// for creating account

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Account
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}
	created := account(user)

	json.NewEncoder(w).Encode(created)
}

func account(user models.Account) models.Account {

	if len(user.BankAccount) > 10 || len(user.IFSC) > 8 {
		log.Fatal("incorrect credentials")
	}

	user.DailyLimit = 1000
	user.CurrentBalance = 0
	user.LastUpdated = time.Now()
	_, err := controller.IdCollection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return user

}

// TODO: for withdrawl of money

func Withdrawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.Withdraw
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	// finding user in database
	var userVerify models.Account
	filter := bson.M{"ifsc": user.IFSC, "bankaccount": user.BankAccount}
	err = controller.IdCollection.FindOne(context.TODO(), filter).Decode(&userVerify)
	if err != nil {
		log.Fatal(err)
	}

	// update daily limit

	now := time.Now() // Get the current time.

	// Check if the last reset was before the start of the current day.
	if userVerify.LastUpdated.Before(now.Truncate(24 * time.Hour)) {
		userVerify.DailyLimit = 1000 // Reset the daily usage to 1000.
		userVerify.LastUpdated = now // Update the last reset time to the current time.
	}

	if userVerify.Status == "INACTIVE" {
		log.Fatal("status inactive")
	}

	if userVerify.DailyLimit < user.Amount {
		log.Fatal("limit exceeded")
	}

	if userVerify.CurrentBalance < user.Amount {
		log.Fatal("invalid amount")
	}

	userVerify.CurrentBalance -= user.Amount
	userVerify.DailyLimit -= user.Amount
	err = controller.IdCollection.FindOneAndReplace(context.TODO(), filter, userVerify).Decode(&userVerify)
	if err != nil {
		log.Fatal(err)
	}
	userVerify.CurrentBalance -= user.Amount
	userVerify.DailyLimit -= user.Amount
	json.NewEncoder(w).Encode(userVerify)

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
	var userVerify models.Account
	filter := bson.M{"ifsc": user.IFSC, "bankaccount": user.BankAccount}
	err = controller.IdCollection.FindOne(context.TODO(), filter).Decode(&userVerify)
	if err != nil {
		log.Fatal("user not found")
	}

	if userVerify.Status == "INACTIVE" {
		log.Fatal("status inactive")
	}

	// update daily limit
	now := time.Now() // Get the current time.

	// Check if the last reset was before the start of the current day.
	if userVerify.LastUpdated.Before(now.Truncate(24 * time.Hour)) {
		userVerify.DailyLimit = 1000 // Reset the daily usage to 1000.
		userVerify.LastUpdated = now // Update the last reset time to the current time.
	}

	if userVerify.DailyLimit < user.Amount {
		log.Fatal("limit exceeded")
	}

	if user.PaymentMode == "debit" && !userVerify.AllowDebit {
		log.Fatal("debit payment mode not allowed for this account")
	}

	if user.PaymentMode == "credit" && !userVerify.AllowCredit {
		log.Fatal("credit payment mode not allowed for this account")
	}

	userVerify.CurrentBalance += user.Amount
	userVerify.DailyLimit -= user.Amount
	err = controller.IdCollection.FindOneAndReplace(context.TODO(), filter, userVerify).Decode(&userVerify)
	if err != nil {
		log.Fatal(err)
	}
	userVerify.CurrentBalance += user.Amount
	userVerify.DailyLimit -= user.Amount
	json.NewEncoder(w).Encode(userVerify)

}

// get balance

func GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// decoding the response
	var user models.Balance
	fmt.Println(user)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	var userVerify models.Account
	filter := bson.M{"ifsc": user.IFSC, "bankaccount": user.BankAccount}
	err = controller.IdCollection.FindOne(context.TODO(), filter).Decode(&userVerify)
	if err != nil {
		log.Fatal("user not found")
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": userVerify.CurrentBalance})
}
