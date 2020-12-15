package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const filename = "assets/authuser/passwd.csv"

func main() {
	fmt.Println("* bsshd user manager")
	fmt.Println("  * 1. add user")
	fmt.Println("  * 2. delete user")
	fmt.Println("  * 3. verify")
	fmt.Println("  * 4. exit")
	for {
		fmt.Print("type number (1-4) > ")
		var s int
		fmt.Scan(&s)

		switch s {
		case 1:
			if ok := adduser(); !ok {
				fmt.Println("add user failed")
			} else {
				fmt.Println("add user success!")
			}
		case 2:
			if ok := deleteuser(); !ok {
				fmt.Println("delete user failed")
			} else {
				fmt.Println("delete user success!")
			}
		case 3:
			if ok := verify(); !ok {
				fmt.Println("verify user failed")
			} else {
				fmt.Println("verify user success!")
			}
		case 4:
			fmt.Println("exit")
			return
		default:
			fmt.Println("invalid number.")
			continue
		}
	}
}

type authuser map[string]string

func adduser() bool {
	au := load(filename)

	var (
		name   string
		passwd string
	)
	fmt.Print("username: ")
	fmt.Scan(&name)
	fmt.Print("password: ")
	fmt.Scan(&passwd)

	if _, ok := au[name]; ok {
		fmt.Printf("user %v is already exist.\n", name)
		return false
	}

	hash, e := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if e != nil {
		panic(e)
	}
	au[name] = string(hash)

	save(filename, au)
	return true
}

func deleteuser() bool {
	au := load(filename)

	var name string
	fmt.Print("username: ")
	fmt.Scan(&name)

	if _, ok := au[name]; !ok {
		fmt.Printf("user: %s is not exist.\n", name)
		return false
	}

	delete(au, name)
	save(filename, au)
	return true
}

func verify() bool {
	au := load(filename)

	var (
		name   string
		passwd string
	)
	fmt.Print("username: ")
	fmt.Scan(&name)
	fmt.Print("password: ")
	fmt.Scan(&passwd)

	if _, ok := au[name]; !ok {
		fmt.Printf("user: %s is not exist.\n", name)
		return false
	}

	hash := []byte(au[name])
	e := bcrypt.CompareHashAndPassword(hash, []byte(passwd))
	if e != nil {
		fmt.Printf("password: %s is invalid\n", passwd)
		return false
	}
	return true
}

func load(filename string) authuser {
	csvBytes, e := ioutil.ReadFile(filename)
	if e != nil {
		panic(e)
	}
	csvContent := string(csvBytes)
	au := make(authuser)
	for _, l := range strings.Split(csvContent, "\n") {
		up := strings.Split(l, ",")
		if len(up) != 2 { // skip if file is empty
			break
		}
		au[up[0]] = up[1]
	}
	return au
}

func save(filename string, au authuser) {
	var csvContent string
	for user, passwdHash := range au {
		csvContent += fmt.Sprintf("%s,%s\n", user, passwdHash)
	}

	csvBytes := []byte(csvContent)

	e := ioutil.WriteFile(filename, csvBytes, os.ModePerm)
	if e != nil {
		panic(e)
	}
}
