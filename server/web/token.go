package web

import (
	"fmt"
	"time"
)

type token struct {
	user   string
	tkn    string
	expire time.Time
}

var (
	tknUser  = map[string]*token{}
	userTkn  = map[string]*token{}
	duration = time.Hour
)

func addToken(u string) (tkn string) {
	now := time.Now()
	if t, ok := userTkn[u]; ok {
		t.expire = now.Add(duration)
		return t.tkn
	}

	tkn = fmt.Sprintf("%d", now.UnixNano())
	t := &token{
		user:   u,
		tkn:    tkn,
		expire: now.Add(duration),
	}
	tknUser[tkn] = t
	userTkn[u] = t
	return
}

func rmTknUser(tkn string) {
	if t, ok := tknUser[tkn]; ok {
		delete(tknUser, tkn)
		delete(userTkn, t.user)
		return
	}
}

func rmUserTkn(u string) {
	if t, ok := userTkn[u]; ok {
		delete(tknUser, t.tkn)
		delete(userTkn, u)
		return
	}
}

func getTknUser(tkn string) (string, bool) {
	t, ok := tknUser[tkn]
	if !ok {
		return "", false
	}
	now := time.Now()
	if now.After(t.expire) {
		rmTknUser(tkn)
		return "", false
	}
	t.expire = now.Add(duration)
	return t.user, true
}

func getUserTkn(u string) (string, bool) {
	t, ok := userTkn[u]
	if !ok {
		return "", false
	}
	now := time.Now()
	if now.After(t.expire) {
		rmUserTkn(u)
		return "", false
	}
	t.expire = now.Add(duration)
	return t.tkn, true
}
