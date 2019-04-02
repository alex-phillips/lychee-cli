package lychee

import "encoding/json"

type Album struct {
	ID           json.Number `json:",Number"`
	Title        string
	Public       string
	Description  string
	Visible      string
	Downloadable string
	License      string
	Sysdate      string
	MinTakestamp string
	MaxTakestamp string
	Password     string
}

func (album *Album) GetID() string {
	return album.ID.String()
}

type Photo struct {
	ID     json.Number
	Public string
	Size   string
	Tags   string
	Title  string
	Type   string
	URL    string
	Width  int
	Height int
}

func (photo *Photo) GetID() string {
	return photo.ID.String()
}

type Search struct {
	Albums []Album
	Hash   string
	Photos []Photo
}

type AlbumsGetResponse struct {
	Albums []Album
}
