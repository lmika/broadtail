package ytrss

import (
	"encoding/xml"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Entries []Entry  `xml:"http://www.w3.org/2005/Atom entry"`
}

type Entry struct {
	VideoID   string    `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`
	ChannelID string    `xml:"http://www.youtube.com/xml/schemas/2015 channelId"`
	Title     string    `xml:"http://www.w3.org/2005/Atom title"`
	Link      string    `xml:"http://www.w3.org/2005/Atom link"`
	Published time.Time `xml:"http://www.w3.org/2005/Atom published"`
}
