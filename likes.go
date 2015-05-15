package vkutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
)

func (api *Api) LikesAdd(likeType string, ownerId, itemId int,
	args ...url.Values) error {
	resp, err := api.vkApi.Request(vk.METHOD_LIKES_ADD, url.Values{
		"type":     {likeType},
		"item_id":  {fmt.Sprint(itemId)},
		"owner_id": {fmt.Sprint(ownerId)},
	})
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	var r ResponseLikes
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return err
	}
	return nil
}

func (api *Api) LikesIsLiked(likeType string, ownerId, itemId int,
	args ...url.Values) (liked bool, copied bool, err error) {
	resp, err := api.vkApi.Request(vk.METHOD_LIKES_IS_LIKED, url.Values{
		"type":     {likeType},
		"item_id":  {fmt.Sprint(itemId)},
		"owner_id": {fmt.Sprint(ownerId)},
	})
	if err != nil {
		return
	}
	type ResponseLikesIsliked struct {
		Response struct {
			Liked  int `json:"liked"`
			Copied int `json:"copied"`
		} `json:"response"`
		ResponseError
	}
	var r ResponseLikesIsliked
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return
	}
	if r.Response.Liked == 1 {
		liked = true
	}
	if r.Response.Copied == 1 {
		copied = true
	}
	if r.Error.Msg != "" {
		err = errors.New(r.Error.Msg)
	}

	return
}
