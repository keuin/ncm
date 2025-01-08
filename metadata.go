package ncm

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Metadata struct {
	MusicID       int64    `json:"musicId"`
	MusicName     string   `json:"musicName"`
	Artist        []Artist `json:"artist"`
	AlbumID       int64    `json:"albumId"`
	Album         string   `json:"album"`
	AlbumPicDocID string   `json:"albumPicDocId"`
	AlbumPic      string   `json:"albumPic"`
	Bitrate       int      `json:"bitrate"`
	Mp3DocID      string   `json:"mp3DocId"`
	Duration      int      `json:"duration"`
	MvID          int64    `json:"mvId"`
	Alias         []string `json:"alias"`
	TransNames    []string `json:"transNames"`
	// Format is music file extension name, e.g. "mp3"
	Format string `json:"format"`
}

type Artist struct {
	ID   int64
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
