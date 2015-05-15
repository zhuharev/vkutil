package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
	"strings"
)

func (api *Api) FriendsGet(args ...url.Values) ([]int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_FRIENDS_GET, args...)
	if err != nil {
		return nil, err
	}
	var r ResponseIds
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	return r.Resp.Items, nil
}

func (api *Api) FriendsGetRequests(args ...url.Values) ([]int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_FRIENDS_GET_REQUESTS, args...)
	if err != nil {
		return nil, err
	}
	var r ResponseIds
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	return r.Resp.Items, nil
}

func (api *Api) FriendsGetAllFollowers() ([]int, error) {
	ids, err := api.FriendsGet()
	if err != nil {
		return nil, err
	}
	followers, err := api.UsersGetFollowers()
	if err != nil {
		return nil, err
	}
	followers = append(followers, ids...)
	return followers, nil
}

type ResponseAreFriends struct {
	Response []AreFriends `json:"response"`
	ResponseError
}

type AreFriends struct {
	Id           int    `json:"user_id"`
	FriendStatus int    `json:"friend_status"`
	Sign         string `json:"sign"`
}

func (api *Api) FriendsAreFriends(ids []int,
	args ...url.Values) ([]AreFriends, error) {
	resp, err := api.vkApi.Request(vk.METHOD_FRIENDS_ARE_FRIENDS, url.Values{
		"user_ids": {strings.Join(arrIntToStr(ids), ",")},
	})
	if err != nil {
		return nil, err
	}
	var r ResponseAreFriends
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}
	return r.Response, nil
}

type ResponseInt struct {
	Response int `json:"response"`
	ResponseError
}

func (api *Api) FriendsAdd(userId int, args ...url.Values) (int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_FRIENDS_ADD, url.Values{
		"user_id": {fmt.Sprint(userId)},
	})
	if err != nil {
		return 0, err
	}
	var r ResponseInt
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return 0, err
	}
	return r.Response, nil
}
