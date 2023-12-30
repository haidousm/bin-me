package main

type contextKey string

func (ck contextKey) String() string {
	return string(ck)
}

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const authenticatedUserIDKey = contextKey("authenticatedUserID")
