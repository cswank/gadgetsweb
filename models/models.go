package models

import (
	"os"
	"labix.org/v2/mgo/bson"
	"io/ioutil"
)

var (
	DBPath = os.Getenv("GADGETSDB")
)

type DB struct {
	Users map[string]User     `bson:"users"`
	Gadgets map[string]Gadget `bson:"gadgets"`
	Methods map[string]Method `bson:"methods"`
}

func getDB() *DB {
	db := &DB{}
	err := db.Open()
	if err != nil {
		db.Gadgets = map[string]Gadget{}
		db.Users =  map[string]User{}
		db.Methods =  map[string]Method{}
	}
	return db
}

type Gadget struct {
	Name string `bson:"name" json:"name"`
	Host string `bson:"host" json:"host"`
}

type User struct {
	Id uint64 `bson:"0" json:"-"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"-" json:"password"`
	HashedPassword []byte `bson:"hashedPassword" json:"-"`
}

type Method struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Gadget string `json:"gadget"`
	Steps []string `json:"steps"`
}

type Timeseries struct {
	Name     string              `bson:"name" json:"name"`
	Data     []interface{}       `bson:"data" json:"data"`
}

type Summary struct {
	Location string `json:"location" json:"location"`
	Name string `json:"name" json:"name"`
	Direction string `json:"direction" json:"direction"`
}

func (d *DB) Open() error {
	b, err := ioutil.ReadFile(DBPath)
	if err == nil {
		err = bson.Unmarshal(b, d)
	}
	return err
}

func (d *DB) Save() error {
	b, err := bson.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(DBPath, b, 0644)
}
