package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	vk "github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) WallGet(ownerId int, filter ...url.Values) ([]Post, error) {
	vals := url.Values{}
	if len(filter) == 1 {
		vals = filter[0]
	}
	vals.Set("owner_id", fmt.Sprint(ownerId))
	resp, err := api.vkApi.Request(vk.METHOD_WALL_GET, vals)
	if err != nil {
		return nil, err
	}
	var r ResponsePosts
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	if r.Error.Msg != "" {
		return nil, errors.New(r.Error.Msg)
	}
	return r.Response.Items, nil
}
