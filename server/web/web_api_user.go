package web

import (
	"log"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*User) Nav(user string) (ret Result) {
	log.Printf("nav user:%s \n", user)
	if user == "admin" {
		ret.Data = defNav
	} else {
		ret.Data = defNav[0]
	}
	return
}
