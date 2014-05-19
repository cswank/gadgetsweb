package models

import (
	"reflect"
)

var (
	getGadgetsQuery = "SELECT name, host FROM gadgets"
	saveGadgetQuery = "INSERT INTO gadgets (name, host) VALUES (?, ?)"
)


func GetGadgets() []Gadget {
	// db, err := getDB()
	// if err != nil {
	// 	return []Gadget{}
	// }

	//gadgets := db.Use("gadgets")
	return  []Gadget{}
}

func (g *Gadget) toMap() map[string]interface{} {
	m := map[string]interface{}{}
	s := reflect.ValueOf(g).Elem()
        typeOfT := s.Type()
        for i := 0; i < s.NumField(); i++ {
                f := s.Field(i)
                m[typeOfT.Field(i).Name] = f.Interface()
        }
	return m
}


func (g *Gadget)Save() error {
	return nil
}

func (g *Gadget)Delete() error {
	return nil
}
