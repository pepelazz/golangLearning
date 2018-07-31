package main

import "fmt"

const (
	ADMIN = "admin"
)

func main() {
	user := getUser()
	fmt.Printf("\n isAdmin: %s\n", user.GetBalance())
}

type User struct {
	Name string
	Role string
}

type BankUser struct {
	Name string
	Role string
	Balance float64
}

type UserInterface interface {
	IsAdmin() bool
	GetBalance() float64
}

func (u *User) IsAdmin() bool  {
	return u.Role == ADMIN
}

func (u *User) GetBalance() float64  {
	return 12.
}

func (u *BankUser) IsAdmin() bool  {
	return u.Role == ADMIN
}

func (u *BankUser) GetBalance() float64  {
	return u.Balance
}

func getUser() UserInterface {
	return &BankUser{Role: ADMIN, Balance: 90.}
}