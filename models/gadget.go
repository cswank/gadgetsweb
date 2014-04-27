package models

var (
	getGadgetsQuery = "SELECT name, host FROM gadgets"
	saveGadgetQuery = "INSERT INTO gadgets (name, host) VALUES (?, ?)"
)


func GetGadgets() []Gadget {
	db := getDB()
	gadgets := make([]Gadget, len(db.Gadgets))
	i := 0
	for _, val := range db.Gadgets {
		gadgets[i] = val
		i += 1
	}
	return gadgets
}

func (g *Gadget)Save() error {
	db := getDB()
	db.Gadgets[g.Name] = *g
	return db.Save()
}

func (g *Gadget)Delete() error {
	db := getDB()
	delete (db.Gadgets, g.Name)
	return db.Save()
}
