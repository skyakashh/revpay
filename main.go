package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skyakashh/revpay/routers"
)

func main() {
	fmt.Println("mongo db api")
	r := routers.Router()
	fmt.Println("server is getting started ...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("listen and serve 4000")
}
