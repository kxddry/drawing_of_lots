package main

import (
	"github.com/joho/godotenv"
	"log"
	"strconv"
)

const (
	MaxUsers = 32
)

var (
	envFile, errOpen     = godotenv.Read("pkg/.env")          // the .env file is located in the pkg folder
	token                = envFile["TOKEN"]                   // the telegram bot API token
	owner                = envFile["OWNER"]                   // the ID of the one starting the polls
	int64owner, errParse = strconv.ParseInt(owner, 10, 64)    // 64 bytes is required for userIDs
	peers                = make([]int64, 0, MaxUsers)         // users
	usersHashmap         = make(map[int64][]string, MaxUsers) // {...chatID: [username, firstName]...}
)

func init() {
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	if token == "" {
		log.Fatal("check your API token")
	}
	if owner == "" {
		log.Fatal("No owner provided")
	}
	if errParse != nil {
		log.Fatal(errParse)
	}
	// if any issues during reading the env file arise, stop the program
}
