package models

import (
	"os"
	"github.com/HouzuoGuo/tiedot/db"
	"math/rand"
	"time"
)

var (
	DBPath = os.Getenv("GADGETSDB")
)

func getDB() (*db.DB, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	return db.OpenDB(DBPath)
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
	Permission string `bson:"permission" json:"permission"`
}

type Method struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Gadget string `json:"gadget"`
	Steps []string `json:"steps"`
}

type Note struct {
	Text string     `json:"name"`
	Name string     `json:"gadget"`
	Taken time.Time `json:"steps"`
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

