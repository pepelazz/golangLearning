package main

type UserAccountInterface interface {
	GetAvatar() string
}

type UserAccountFb struct {
	Avatar string `json:"avatar"`
}

type UserAccountVk struct {
	Photo string `json:"photo"`
}

func (a *UserAccountFb) GetAvatar() string {
	return a.Avatar
}

func (a *UserAccountVk) GetAvatar() string {
	return a.Photo
}
