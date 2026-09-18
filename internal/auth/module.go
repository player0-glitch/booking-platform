package auth

import "github.com/gorilla/sessions"

type AuthModule struct {
	session sessions.Session
}

func (a *AuthModule) NewAuthModule(s sessions.Session) *AuthModule {
	return &AuthModule{
		session: s,
	}
}
