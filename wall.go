package vkutil

import (
	"encoding/json"
	"fmt"
	vk "github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) WallGet(ownerId int, filter ...Opts) ([]Post, error) {
	vals := url.Values{
		"owner_id": {fmt.Sprint(ownerId)},
	}
	if len(filter) > 0 && filter[0].Filter != "" {
		vals.Set("filter", filter[0].Filter)
	}
	resp, err := api.vkApi.Request(vk.METHOD_WALL_GET, vals)
	if err != nil {
		return nil, nil
	}
	var r ResponsePosts
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, nil
	}
	return nil, nil
}
