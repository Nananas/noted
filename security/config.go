package security

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func LoadConfig() {
	cache = parseAccountsFile()
}

func parseAccountsFile() map[string]string {
	bytes, err := ioutil.ReadFile(ACCOUNTS_FILEPATH)
	if err != nil {
		log.Fatal(err)
	}

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
			hash := split[1]

			fmt.Println("Accounts found: ")
			fmt.Println("> " + user + " - " + hash)

			c[user] = hash

		}

	}

	return c

}
