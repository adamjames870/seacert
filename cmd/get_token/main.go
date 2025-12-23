package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func main() {

	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}

	url := os.Getenv("SUPABASE_URL")
	apiKey := os.Getenv("SUPABASE_PERISHABLE_KEY") //("SUPABASE_ANON_KEY")
	password := os.Getenv("TEST_USER_PASSWORD")

	var email string
	if len(os.Args) > 1 {
		email = os.Args[1]
	} else {
		email = os.Getenv("TEST_USER_EMAIL")
	}

	if email == "" {
		fmt.Println("Usage: go run ./cmd/get_token <email>")
		fmt.Println("Or set TEST_USER_EMAIL environment variable")
		return
	}

	client, errClient := supabase.NewClient(url, apiKey, nil)
	if errClient != nil {
		panic(errClient)
	}

	authResp, errResp := client.Auth.SignInWithEmailPassword(email, password)
	if errResp != nil {
		panic(errResp)
	}

	fmt.Println("Access token:", authResp.Session.AccessToken)
	fmt.Println("Refresh token:", authResp.Session.RefreshToken)
}
