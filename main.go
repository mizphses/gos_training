package main

import(
	"bytes"
	"golang.org/x/crypto/argon2"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type userAccount struct {
	Username string `json:"id"`
	Password string `json:"pwd"`
}

var userAccounts = []userAccount{
	{Username: "user1", Password: "pwd1"},
	{Username: "user2", Password: "pwd2"},
}

func main() {
	http.HandleFunc("/get", uketori)
	http.HandleFunc("/post", haitatsu)
	http.ListenAndServe(":8080", nil)
}

func uketori(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if accountInfo := enc.Encode(&userAccounts); accountInfo != nil {
		log.Print(accountInfo)
	}
	fmt.Println(buf.String())


	_, accountInfo := w.Write(buf.Bytes())
	if accountInfo != nil {
		return
	}
}

func haitatsu(w http.ResponseWriter, r *http.Request) {
	var accountInfoB userAccount
	json.NewDecoder(r.Body).Decode(&accountInfoB)
	pwd_raw := []byte(accountInfoB.Password)
	accountInfoB.Password = string(argon2.IDKey(pwd_raw, []byte("salt"), 1, 64*1024, 4, 32))
	fmt.Fprintf(w, "%sは%sのパスワードをもつアカウントとして登録されました", accountInfoB.Username, accountInfoB.Password)
}

