// Package auth は，bsshdの認証に関わる機能を提供する
package auth

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hiragi-gkuth/bsshd/pkg/ids"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

// PasswordAuthenticationError は，パスワード認証が失敗したときに投げられるエラー
type PasswordAuthenticationError struct {
	password string
	user     string
	message  string
}

const filename = "assets/authuser/passwd.csv"

// Password はパスワード認証時に呼び出され，認証を行う関数
func Password(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	// store password
	authSession := ids.GetAuthSession()
	key := conn.RemoteAddr().String()
	authInfo, _ := authSession.Get(key)
	// use fake password for operation server
	authInfo.Passwords = append(authInfo.Passwords, "")

	// authentication
	users := fetchUserList()
	passwdHash, ok := users[conn.User()]
	if !ok { // user not exists
		goto failure
	}
	if !verify(passwdHash, password) { // incorrect password
		goto failure
	}
	goto success

success: // if authentication success
	authInfo.Results = append(authInfo.Results, "Success")
	authSession.Set(key, authInfo)
	// return &ssh.Permissions{}, nil

	// force failure
	return nil, PasswordAuthenticationError{}

failure: // if authentication failed
	authInfo.Results = append(authInfo.Results, "Failure")
	authSession.Set(key, authInfo)
	return nil, PasswordAuthenticationError{
		password: string(password),
		user:     conn.User(),
	}
}

// PasswordHoneyPot は，ハニーポットサーバ動作時に呼び出される
func PasswordHoneyPot(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	// store password
	authSession := ids.GetAuthSession()
	key := conn.RemoteAddr().String()
	authInfo, _ := authSession.Get(key)
	authInfo.Passwords = append(authInfo.Passwords, string(password))
	authInfo.Results = append(authInfo.Results, "Failure")
	authSession.Set(key, authInfo)

	// force failure
	return nil, PasswordAuthenticationError{}
}

type authuser map[string]string

func fetchUserList() authuser {
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

func verify(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	e := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	return e == nil
}

func (pae PasswordAuthenticationError) Error() string {
	return fmt.Sprintf("Failed login for %v using %v", pae.user, pae.password)
}
