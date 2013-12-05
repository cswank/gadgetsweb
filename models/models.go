package models

import (
	"time"
)


type Device struct {
	Units string      `json:"units"`
	Value interface{} `json:"value"`
	ID    string      `json:"id"`
}

type Location struct {
	Input  map[string]Device `json:"input"`
	Output map[string]Device `json:"output"`
}

type Gadgets struct {
	Sender      string      `json:"sender"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Locations   map[string]Location    `json:"locations"`
}

type Timeseries struct {
	Name     string              `json:"name"`
	Data     []interface{}       `json:"data"`
}

type Summary struct {
	Location string `json:"location"`
	Name string `json:"name"`
	Direction string `json:"direction"`
}
