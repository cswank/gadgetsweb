package models

import (
	"log"
	"time"
)

var (
	saveNoteQuery = "INSERT INTO notes (text, gadget, taken) VALUES (?, ?, ?)"
	getNotesQuery = "select text, gadget, taken from notes where gadget = ? and taken >= ? and taken <= ? ORDER BY taken DESC"
)

type Note struct {
	Text   string    `json:"text"`
	Gadget string    `json:"gadget"`
	Taken  time.Time `json:"time"`
}

func GetNotes(gadget string, start, end time.Time) []Note {
	notes := []Note{}
	rows, err := DB.Query(getNotesQuery, gadget, start.Unix(), end.Unix())
	if err != nil {
		log.Println(err)
		return []Note{}
	}
	for rows.Next() {
		var ts int64
		n := Note{}
		err = rows.Scan(&n.Text, &n.Gadget, &ts)
		if err != nil {
			return []Note{}
		}
		n.Taken = time.Unix(ts, 0)
		notes = append(notes, n)
	}
	return notes
}

// func (n *Note)Delete() error {
// 	return nil
// }

func (n *Note) Save() error {
	if n.Taken.Equal(time.Time{}) {
		n.Taken = time.Now()
	}
	_, err := DB.Query(saveNoteQuery, n.Text, n.Gadget, n.Taken)
	return err
}
