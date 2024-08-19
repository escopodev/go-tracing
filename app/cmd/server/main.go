package main

import (
	"log"

	checkouts "github.com/escopodev/checkouts/cmd/server"
	payments "github.com/escopodev/payments/cmd/server"
)

func main() {
	log.Println("running applications")

	go func() {
		checkouts.Serve()
	}()

	payments.Serve()
}
