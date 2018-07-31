package main

import "encoding/gob"

const (
	HTTP_PORT = 3010
)

var (
	userList []User
)

func main() {
	gob.Register(&UserAccountVk{})
	gob.Register(&UserAccountFb{})
	// создаем список пользователей
	initUserList()
	conncetDb()
	startWebServer()
}

func initUserList() {
	userList = []User{{Id: 1, Name: "Иванов", Social: "fb"}, {Id: 2, Name: "Сидоров", Social: "vk"}}
}
