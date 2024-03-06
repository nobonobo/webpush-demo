//go:build ignore

package main

import (
	"log"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	// Generate vapid keys
	vapidPrivateKey, vapidPublicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("VAPID private key:", vapidPrivateKey)
	log.Println("VAPID public key:", vapidPublicKey)
}
