package entities

import "time"

type Project struct {
	Id           int
	Title        string
	Sdate        time.Time
	Edate        time.Time
	Duration     string
	Content      string
	Technologies []string
	Tnode        bool
	Treact       bool
	Tjs          bool
	Thtml        bool
}

var Data = map[string]interface{}{
	"Title":     "Personal Web",
	"IsLogin":   true,
	"Id":        1,
	"UserName":  "Tommy",
	"FlashData": "",
}
