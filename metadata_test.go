package ncm

import (
	"encoding/json"
	"testing"
)

func TestArtist_UnmarshalJSON(t *testing.T) {
	type fields struct {
		ID   int64
		Name string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				ID:   123,
				Name: "name",
			},
			args: args{
				data: []byte(`["name",123]`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Artist{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := a.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadata_UnmarshalJSON(t *testing.T) {
	m := new(Metadata)
	err := json.Unmarshal([]byte(`{
    "musicId": 123,
    "musicName": "abc",
    "artist": [
        [
            "abc",
            123
        ]
    ],
    "albumId": 123,
    "album": "abc",
    "albumPicDocId": "123123",
    "albumPic": "abc",
    "bitrate": 123123,
    "mp3DocId": "abc",
    "duration": 123,
    "mvId": 123,
    "alias": [
        "abc"
    ],
    "transNames": [
        "abc"
    ],
    "format": "mp3"
}`), m)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Artist) != 1 {
		t.Fail()
	}
	if m.MusicID != 123 {
		t.Fail()
	}
}
