package main

import (
	"strconv"
	"encoding/json"
	"encoding/gob"
	"bytes"
)

const (
	ENCODE_TYPE_JSON = "json"
	ENCODE_TYPE_GOB  = "gob"
)

type User struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Social      string `json:"social"`
	UserAccount UserAccountInterface `json:"-"`
}

func (u *User) SaveToDB() (error) {
	switch u.Social {
	case "vk":
		u.UserAccount = &UserAccountVk{}
	case "fb":
		u.UserAccount = &UserAccountFb{}

	}
	return localDb.PutGob(USER_DB_BUCKET, strconv.Itoa(u.Id), u)
}

func getUserListFromDB(encodeType string) (res []User, err error) {
	byteMap := map[string][]byte{}
	localDb.GetBucketList(USER_DB_BUCKET, byteMap)

	for _, v := range byteMap {
		user := User{}
		switch encodeType {
		case ENCODE_TYPE_JSON:
			err := json.Unmarshal(v, &user)
			if err != nil {
				return nil, err
			}
		case ENCODE_TYPE_GOB:
			buf := bytes.NewBuffer(v)
			dec := gob.NewDecoder(buf)
			err = dec.Decode(&user)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, user)
	}

	return
}
