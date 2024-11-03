package main

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("social_security_number ", "123-45-6789")
	log.Info().Msg("email ", "testing@test.com")
	log.Printf("User email: foobar")
	fmt.Println("Customer phone number: 123-456-7890")
}
