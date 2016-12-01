package main

import (
	"fmt"

	"github.com/asiainfoLDP/datafoundry_payment/api/openshift"
)

func demo() {
	openshift.Init()
}

func init() {
	fmt.Println("TEST DEMO....STARTED")
	demo()
	fmt.Println("TEST DEMO....FINISHED.")
}
