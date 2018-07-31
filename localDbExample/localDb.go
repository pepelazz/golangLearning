package main

import (
	"os"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
	"encoding/json"
	"encoding/gob"
	"bytes"
	"log"
)

const (
	USER_DB_BUCKET = "user"
)

var (
	localDb *LocalDbType
	//boltDb  *bolt.DB
)

type LocalDbType struct{
	boltDb  *bolt.DB
}

func conncetDb() {
	// boltdb
	localDb = &LocalDbType{}
	err := localDb.Open()
	checkErr(err, "Connect to boltDb")
}

func (bd *LocalDbType) Open() (err error) {
	err = os.Mkdir("boltDb", os.ModePerm)
	if err != nil {
		println(fmt.Sprintf("LocalDbType.Open: %s", err))

	}
	bd.boltDb, err = bolt.Open("boltDb/my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

func (bd *LocalDbType) Close() {
	bd.boltDb.Close()
}


func (bd *LocalDbType) GetJson(bName, key string, data interface{}) error {
	return bd.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bName))
		if bucket != nil {
			v := bucket.Get([]byte(key))
			if len(v) > 0 {
				err := json.Unmarshal(v, &data)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (bd *LocalDbType) GetStr(bName, key string) (data string, err error) {
	err = bd.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bName))
		if bucket != nil {
			v := bucket.Get([]byte(key))
			data = string(v)
		}
		return nil
	})
	return
}

func (bd *LocalDbType) PutByte(bName, key string, data []byte) (err error) {
	return bd.boltDb.Update(func(tx *bolt.Tx) error {
		// создаем bucket с нужным именем если он еще не создан
		b, err := tx.CreateBucketIfNotExists([]byte(bName))
		if err != nil {
			return err
		}

		return b.Put([]byte(key), data)
	})
}

func (bd *LocalDbType) PutJson(bName, key string, data interface{}) error {
	// сериализауем в json
	encoded, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return bd.PutByte(bName, key, encoded)
}

func (bd *LocalDbType) PutGob(bName, key string, data interface{}) (err error) {
	// сериализуем в gob
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(data)
	if err != nil {
		return
	}

	return bd.PutByte(bName, key, buf.Bytes())
}

func (bd *LocalDbType) GetGob(bName, key string, data interface{}) error {
	return bd.boltDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bName))
		if bucket != nil {
			v := bucket.Get([]byte(key))
			if len(v) > 0 {
				buf := bytes.NewBuffer(v)
				dec := gob.NewDecoder(buf)
				return dec.Decode(data)
			}
		}
		return nil
	})
}

func (bd *LocalDbType) Delete(bName, key string) error {
	return bd.boltDb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bName))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

func (bd *LocalDbType) GetBucketList(bName string, res map[string][]byte) error {
	return bd.boltDb.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bName))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res[string(k)] = v
		}
		return nil
	})
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
