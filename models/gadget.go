package models

var (
	getGadgetsQuery = "SELECT name, host FROM gadgets"
	saveGadgetQuery = "INSERT INTO gadgets (name, host) VALUES (?, ?)"
	deleteGadgetQuery = "DELETE FROM gadgets where name = ?"
)

type GadgetHosts struct {
	Gadgets []Gadget `json:"gadgets"`
}

type Gadget struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

func GetGadgets() (*GadgetHosts, error) {
	db, err := getDB()
	defer db.Close()
	gadgets := &GadgetHosts{}
	rows, err := db.Query(getGadgetsQuery)
	if err != nil {
		return gadgets, err
	}
	for rows.Next() {
		g := Gadget{}
		if err = rows.Scan(&g.Name, &g.Host); err != nil {
			return gadgets, err
		}
		gadgets.Gadgets = append(gadgets.Gadgets, g)
	}
	return gadgets, nil
}

func (g *Gadget)Save() error {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Query(saveGadgetQuery, g.Name, g.Host)
	return err
}

func (g *Gadget)Delete() error {
	db, err := getDB()
	defer db.Close()
	if err != nil {
		return err
	}
	_, err = db.Query(deleteGadgetQuery, g.Name)
	return err
}
