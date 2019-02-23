package main

import "fmt"

const (
	// Prefix for keys
	Prefix = "sylsess_"
	// UserEmailKey key
	UserEmailKey = Prefix + `u:%v:email`
	// UserIDKey key
	UserIDKey = Prefix + `u:%vid`
	// UserLastActiveKey key
	UserLastActiveKey = Prefix + `u:%vts`
)

// KeyUserEmail returns the user email key
func KeyUserEmail(id string) string {
	return fmt.Sprintf(UserEmailKey, id)
}

// KeyUserID returns the user id key
func KeyUserID(id string) string {
	return fmt.Sprintf(UserIDKey, id)
}

// KeyUserLastActive returns the last active key
func KeyUserLastActive(id string) string {
	return fmt.Sprintf(UserLastActiveKey, id)
}
