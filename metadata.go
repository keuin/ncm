package ncm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type Metadata struct {
	MusicID       ID       `json:"musicId"`
	MusicName     string   `json:"musicName"`
	Artist        []Artist `json:"artist"`
	AlbumID       ID       `json:"albumId"`
	Album         string   `json:"album"`
	AlbumPicDocID string   `json:"albumPicDocId"`
	AlbumPic      string   `json:"albumPic"`
	Bitrate       int      `json:"bitrate"`
	Mp3DocID      string   `json:"mp3DocId"`
	Duration      int      `json:"duration"`
	MvID          ID       `json:"mvId"`
	Alias         []string `json:"alias"`
	TransNames    []string `json:"transNames"`
	// Format is music file extension name, e.g. "mp3"
	Format      string    `json:"format"`
	Fee         int       `json:"fee"`
	VolumeDelta float64   `json:"volumeDelta"`
	Privilege   Privilege `json:"privilege"`
}

type Artist struct {
	ID   ID
	Name string
}

func (a *Artist) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	token, err := dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('[') {
		return fmt.Errorf("expected '[' got '%s'", token)
	}
	err = dec.Decode(&a.Name)
	if err != nil {
		return err
	}
	err = dec.Decode(&a.ID)
	if err != nil {
		return err
	}
	token, err = dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim(']') {
		return fmt.Errorf("expected ']' got '%s'", token)
	}
	return nil
}

type ID int64

func (i *ID) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	token, err := dec.Token()
	if err != nil {
		return err
	}
	var value int64
	switch token := token.(type) {
	case string:
		value, err = strconv.ParseInt(token, 10, 64)
	case json.Number:
		value, err = token.Int64()
	default:
		return fmt.Errorf("invalid type of Artist.ID: %v", reflect.TypeOf(token))
	}
	if err != nil {
		return err
	}
	*i = ID(value)
	return nil
}

type Privilege struct {
	Flag int64 `json:"flag"`
}
