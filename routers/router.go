package routers

import (
	"github.com/gorilla/mux"
	"github.com/skyakashh/revpay/mongohelpers"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", mongohelpers.CreateUser).Methods("POST")
	r.HandleFunc("/account", mongohelpers.CreateAccount).Methods("POST")
	r.HandleFunc("/withdraw", mongohelpers.Withdrawl).Methods("PUT")
	r.HandleFunc("/deposit", mongohelpers.Deposit).Methods("PUT")
	return r
}
