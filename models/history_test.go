package models

import (
	"testing"
	"time"
	"fmt"
)

func TestGetHistory(t *testing.T) {
	
	hq := &HistoryQuery{
		Host: "localhost",
		DBName: "gadgets",
		Collection: "updates",
		Location: "hlt",
		Name: "temperature",
		Start: time.Date(2000, 1, 0, 0, 0, 0, 0, time.UTC),
		End: time.Now(),
	}
	h, err := GetHistory(hq)
	
	if err != nil {
		t.Fatal(err)
	}
	s := h[0]
	if len(s.Data) != 18 {
		t.Error(len(s.Data))
	}
	if s.Name != "hlt temperature" {
		t.Error(s.Name)
	}
}


func TestGetDevices(t *testing.T) {
	hq := &HistoryQuery{
		Host: "localhost",
		DBName: "gadgets",
		Collection: "updates",
	}
	d, err := GetDevices(hq)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Links) != 4 {
		t.Error(d)
	}
	fmt.Println(d)
}
