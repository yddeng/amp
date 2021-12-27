package web

import (
	"log"
)

var (
	users = map[string]*User{}
)

func getUser(username string) (u *User, ok bool) {
	u, ok = users[username]
	return
}

type User struct {
	Name     string              `json:"name"`
	Username string              `json:"username,omitempty"`
	Password string              `json:"password,omitempty"`
	Avatar   string              `json:"avatar"`
	Routes   map[string]struct{} `json:"routes"`
}

func (*User) Info(user string) (ret Result) {
	log.Printf("info user:%s \n", user)
	u, _ := getUser(user)
	ret.Data = struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}{
		Name:     u.Name,
		Username: u.Username,
		Avatar:   u.Avatar,
	}
	return
}

func (*User) Nav(user string) (ret Result) {
	log.Printf("nav user:%s \n", user)
	if user == "admin" {
		ret.Data = defNav
		return
	}
	u, _ := getUser(user)
	ret.Data = findNav(u.Routes)
	return
}
