package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nananas/noted/security"
)

var SALT string

func main() {
	log.SetFlags(log.Lshortfile)

	security.SetSalt(SALT)

	bytes, err := ioutil.ReadFile("./.accounts")
	if err != nil {
		log.Fatal(err)
	}

	users := []string{}
	c := make(map[string]string)

	text := string(bytes)

	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		split := strings.Split(lines[i], "|")

		if len(split) < 2 {
			// log.Println("some lines in accounts file are bad:")
			// log.Println(i, lines[i]+".")

		} else {
			user := split[0]
			pass := split[1]

			hash := security.SaltHash(pass)

			// log.Println(user + " - " + pass + " > " + hash)

			users = append(users, user+"|")
			c[user] = hash

		}
	}

	out := ""
	for u, h := range c {
		out = out + u + "|" + h + "\n"
	}

	err = ioutil.WriteFile("./accounts", []byte(out), 0664)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Accounts written in ./accounts file.")

	err = ioutil.WriteFile("./users", []byte(strings.Join(users, "\n")), 0664)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Users written in ./users file.")

	fmt.Println("Everything went well. Should I now remove the .accounts file?\nThis will be better for security, but you will have to recreate this file\neach time you want to recreate the accounts...")
	fmt.Print("[y,yes]: ")

	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	if response == "y" || response == "yes" {

		err = os.Remove("./.accounts")
		if err != nil {
			log.Println(err)
		}

		fmt.Println("Removed .accounts file. Everything went well.")

		return

	}

	fmt.Println("The .accounts file was NOT removed.")

}
