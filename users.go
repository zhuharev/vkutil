package vkutil

import (
	"encoding/json"
	"fmt"
	"github.com/zhuharev/vk"
	"net/url"
	"strings"
)

type ResponseUser struct {
	Response []User `json:"response"`
	ResponseError
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	PhotoId string `json:"photo_id"`
}

func (api *Api) UsersGetFollowers(args ...url.Values) ([]int, error) {
	resp, err := api.vkApi.Request(vk.METHOD_USERS_GET_FOLLOWERS, args...)
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

// utils

func (api *Api) UtilsGetProfilePhoto(ids ...int) ([]User, error) {
	resp, err := api.vkApi.Request(vk.METHOD_USERS_GET, url.Values{
		"user_ids": {strings.Join(arrIntToStr(ids), ",")},
		"fields":   {"photo_id"},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(ids)
	var r ResponseUser
	err = json.Unmarshal(resp, &r)
	return r.Response, err
}

func arrIntToStr(arr []int) (sarr []string) {
	for _, v := range arr {
		sarr = append(sarr, fmt.Sprint(v))
	}
	return
}
