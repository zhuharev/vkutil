package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
)

type PhotoType string

const (
	PHOTO_WALL    = "wall"
	PHOTO_PROFILE = "profile"
	PHOTO_SAVED   = "saved"
)

type Photo struct {
	Id    int `json:"id"`
	Likes struct {
		Count int `json:"count"`
	} `json:"likes"`

	AlbumId int       `json:"album_id"`
	OwnerId int       `json:"owner_id"`
	UserId  int       `json:"user_id"`
	Text    string    `json:"text"`
	Date    EpochTime `json:"date"`
}

func (api *Api) PhotosGet(ownerId int, albumId string, params ...url.Values) ([]Photo, error) {
	rparams := url.Values{}
	rparams.Set("rev", "1")
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("owner_id", fmt.Sprint(ownerId))
	rparams.Set("album_id", albumId)
	resp, e := api.VkApi.Request(vk.METHOD_PHOTOS_GET, rparams)
	if e != nil {
		return nil, e
	}
	var r ResponsePhotos
	e = json.Unmarshal(resp, &r)
	return r.Response.Items, e
}
