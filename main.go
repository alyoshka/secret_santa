package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"time"
)

const body = `
Hello there

I need your help.
Make a present for %s please

Yours faithfully
Santa
`

// host and port of SMTP server
const (
	host = "smtp.gmail.com"
	port = "587"
)

// credentials
var (
	email    string
	password string
)

type person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func init() {
	flag.StringVar(&email, "email", "", "Santa's email address")
	flag.StringVar(&password, "pass", "", "Santa's email password")
}

func main() {
	flag.Parse()
	if email == "" || password == "" {
		fmt.Println("email or password is empty")
		return
	}
	file, err := ioutil.ReadFile("./santas.json")
	if err != nil {
		fmt.Printf("filed to read file, error: %v\n", err)
		return
	}

	persons := []person{}

	err = json.Unmarshal(file, &persons)
	if err != nil {
		fmt.Printf("failed to parse, error: %v\n", err)
		return
	}

	shuffle(persons)

	for i := range persons {
		j := i + 1
		if j == len(persons) {
			j = 0
		}
		err = send(persons[i].Email, persons[j].Name)
		if err != nil {
			fmt.Println("failed to send, error: %s", err)
		}
	}
	fmt.Println("complete")
}

func shuffle(s []person) {
	rand.Seed(time.Now().UnixNano())
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

func send(to string, name string) error {
	msg := "From: Santa <" + email + ">\n" +
		"To: " + to + "\n" +
		"Subject: I need your help\n\n" +
		fmt.Sprintf(body, name)

	auth := smtp.PlainAuth("", email, password, host)
	return smtp.SendMail(host+":"+port, auth, email, []string{to}, []byte(msg))
}
