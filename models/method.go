package models


type Methods struct {
	Methods []Method `json:"methods"`
}

func GetMethods(gadget string) Methods {
	db := getDB()
	m := db.Methods
	methods := Methods{
		Methods: make([]Method, len(m)),
	}
	i := 0
	for _, val := range m {
		methods.Methods[i] = val
		i += 1
	}
	return methods
}

func (m *Method)Save() error {
	db := getDB()
	db.Methods[m.Name] = *m
	return db.Save()
}
