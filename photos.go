package vkutil

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zhuharev/vk"
)

type PhotoType string

const (
	PHOTO_WALL    = "wall"
	PHOTO_PROFILE = "profile"
	PHOTO_SAVED   = "saved"
)

type Size struct {
	Source string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
}

type Photo struct {
	Id    int `json:"id"`
	Likes struct {
		Count int `json:"count"`
	} `json:"likes"`

	Sizes []Size `json:"sizes"`

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

// PhotosGetAll run "photos.getAll" method
func (api *Api) PhotosGetAll(ownerID int, params ...url.Values) ([]Photo, int, error) {
	rparams := url.Values{}
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("owner_id", fmt.Sprint(ownerID))
	rparams.Set("extended", "1")
	rparams.Set("photo_sizes", "1")
	resp, e := api.VkApi.Request(vk.METHOD_PHOTOS_GET_ALL, rparams)
	if e != nil {
		return nil, 0, e
	}
	var r ResponsePhotos
	e = json.Unmarshal(resp, &r)
	return r.Response.Items, r.Response.Count, e
}
