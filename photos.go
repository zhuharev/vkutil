package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
)

type Photo struct {
	Id    int `json:"id"`
	Likes struct {
		Count int `json:"count"`
	} `json:"likes"`
}

func (api *Api) PhotosGet(ownerId int, albumId string, params ...url.Values) ([]Photo, error) {
	rparams := url.Values{}
	rparams.Set("rev", "1")
	if len(params) == 1 {
		rparams = params[0]
	}
	rparams.Set("owner_id", fmt.Sprint(ownerId))
	rparams.Set("album_id", albumId)
	resp, e := api.vkApi.Request(vk.METHOD_PHOTOS_GET, rparams)
	if e != nil {
		return nil, e
	}
	var r ResponsePhotos
	e = json.Unmarshal(resp, &r)
	return r.Response.Items, e
}
